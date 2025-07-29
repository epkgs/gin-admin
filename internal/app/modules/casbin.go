package modules

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gin-admin/internal/configs"
	"gin-admin/internal/defines"
	"gin-admin/internal/dtos"
	"gin-admin/internal/errorx"
	"gin-admin/internal/models"
	"gin-admin/internal/repositories"
	"gin-admin/internal/services"
	"gin-admin/internal/types"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/gormx"
	"gin-admin/pkg/logging"

	"github.com/casbin/casbin/v2"
)

// Load rbac permissions to casbin
type Casbinx struct {
	enforcer *atomic.Value
	ticker   *time.Ticker
	Cache    cachex.Cacher
	// MenuRepo *repositories.Menu
	RoleRepo *repositories.Role
	MenuSvc  *services.Menu
	RoleSvc  *services.Role
}

var _ types.Casbinx = (*Casbinx)(nil)

func InitCasbinx(ctx context.Context, app types.AppContext) (types.Casbinx, error) {
	cb := &Casbinx{
		Cache: app.Cacher(),
		// MenuRepo: repositories.NewMenu(app.DB()),
		RoleRepo: repositories.NewRole(app.DB()),
		MenuSvc:  services.NewMenu(app),
		RoleSvc:  services.NewRole(app),
	}

	app.AddCleaner(ctx, func() {
		cb.Release(context.Background())
	})

	return cb, nil
}

func (a *Casbinx) GetEnforcer() *casbin.Enforcer {
	if v := a.enforcer.Load(); v != nil {
		return v.(*casbin.Enforcer)
	}
	return nil
}

type policyQueueItem struct {
	RoleID string
	Menus  models.Menus
}

func (a *Casbinx) Load(ctx context.Context) error {
	if configs.C.Middleware.Casbin.Disable {
		return nil
	}

	a.enforcer = new(atomic.Value)
	if err := a.load(ctx); err != nil {
		return err
	}

	go a.autoLoad(ctx)
	return nil
}

func (a *Casbinx) load(ctx context.Context) error {
	start := time.Now()
	roles, err := a.RoleRepo.Find(ctx, gormx.WithWhere("status = ?", models.RoleStatus_Enabled), gormx.WithSelect("id"))
	if err != nil {
		return errorx.WrapGormError(ctx, err)
	}

	if len(roles) == 0 {
		return nil
	}

	var resCount int32
	queue := make(chan *policyQueueItem, len(roles))
	threadNum := configs.C.Middleware.Casbin.LoadThread
	lock := new(sync.Mutex)
	buf := new(bytes.Buffer)

	wg := new(sync.WaitGroup)
	wg.Add(threadNum)
	for i := 0; i < threadNum; i++ {
		go func() {
			defer wg.Done()
			ibuf := new(bytes.Buffer)
			for item := range queue {
				for _, res := range item.Menus {
					_, _ = ibuf.WriteString(fmt.Sprintf("p, %s, %s, %s \n", item.RoleID, res.Path, res.Method))
				}
			}
			lock.Lock()
			_, _ = buf.Write(ibuf.Bytes())
			lock.Unlock()
		}()
	}

	for _, item := range roles {
		list, err := a.MenuSvc.List(ctx, dtos.MenuListReq{
			RoleID: item.ID,
			Type:   models.MenuType_BUTTON,
			Pager: dtos.Pager{
				Page: -1,
			},
		})
		if err != nil {
			logging.Error(ctx, "Failed to query role menus", err)
			continue
		}
		atomic.AddInt32(&resCount, int32(len(list.Items)))
		queue <- &policyQueueItem{
			RoleID: item.ID,
			Menus:  list.Items,
		}
	}
	close(queue)
	wg.Wait()

	if buf.Len() > 0 {
		policyFile := configs.C.Middleware.Casbin.GenPolicyFile
		_ = os.Rename(policyFile, policyFile+".bak")
		_ = os.MkdirAll(filepath.Dir(policyFile), 0755)
		if err := os.WriteFile(policyFile, buf.Bytes(), 0666); err != nil {
			logging.Error(ctx, "Failed to write policy file", err)
			return err
		}
		// set readonly
		_ = os.Chmod(policyFile, 0444)

		modelFile := configs.C.Middleware.Casbin.ModelFile
		e, err := casbin.NewEnforcer(modelFile, policyFile)
		if err != nil {
			logging.Error(ctx, "Failed to create casbin enforcer", err)
			return err
		}
		e.EnableLog(configs.C.IsDebug())
		a.enforcer.Store(e)
	}

	logging.Info(ctx, "Casbin load policy",
		map[string]any{
			"cost":      time.Since(start),
			"roles":     len(roles),
			"resources": resCount,
			"bytes":     buf.Len(),
		},
	)
	return nil
}

func (a *Casbinx) autoLoad(ctx context.Context) {
	var lastUpdated int64
	a.ticker = time.NewTicker(time.Duration(configs.C.Middleware.Casbin.AutoLoadInterval) * time.Second)
	for range a.ticker.C {
		val, ok, err := a.Cache.Get(ctx, defines.CacheNSForRole, defines.CacheKeyForSyncToCasbin)
		if err != nil {
			logging.Error(ctx, "Failed to get cache", err, map[string]any{"key": defines.CacheKeyForSyncToCasbin})
			continue
		} else if !ok {
			continue
		}

		updated, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			logging.Error(ctx, "Failed to parse cache value", err, map[string]any{"val": val})
			continue
		}

		if lastUpdated < updated {
			if err := a.load(ctx); err != nil {
				logging.Error(ctx, "Failed to load casbin policy", err)
			} else {
				lastUpdated = updated
			}
		}
	}
}

func (a *Casbinx) Release(ctx context.Context) error {
	if a.ticker != nil {
		a.ticker.Stop()
	}
	return nil
}
