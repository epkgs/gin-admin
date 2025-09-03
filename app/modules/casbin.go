package modules

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"gin-admin/errorx"
	"gin-admin/model/bo"
	"gin-admin/model/dto"
	"gin-admin/model/po"
	"gin-admin/pkg/cachex"
	"gin-admin/pkg/logger"
	"gin-admin/service"
	"gin-admin/types"

	"github.com/casbin/casbin/v2"
)

// Load rbac permissions to casbin
type Casbinx struct {
	app      types.AppContext
	enforcer *atomic.Value
	cacher   cachex.Cacher

	menuSVC *service.Menu
	roleSVC *service.Role

	q *bo.Query

	ticker *time.Ticker
}

var _ types.Casbinx = (*Casbinx)(nil)

func NewCasbinx(ctx context.Context, app types.AppContext) (types.Casbinx, error) {
	cb := &Casbinx{
		app:      app,
		enforcer: new(atomic.Value),
		cacher:   app.Cacher(),
		menuSVC:  service.NewMenu(app),
		roleSVC:  service.NewRole(app),
		q:        bo.Use(app.DB()),
	}

	app.AddCleaner(ctx, func() {
		cb.Release()
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
	Menus  po.Menus
}

func (a *Casbinx) Load(ctx context.Context) error {
	if a.app.Config().Casbin.Disable {
		return nil
	}

	if err := a.load(ctx); err != nil {
		return err
	}

	go a.autoLoad(ctx)
	return nil
}

func (a *Casbinx) load(ctx context.Context) error {
	start := time.Now()

	r := a.q.Role

	roles, err := r.WithContext(ctx).Select(r.ID).Where(r.Status.Eq(po.RoleStatus_Enabled)).Find()
	if err != nil {
		return errorx.WrapGormError(err)
	}

	if len(roles) == 0 {
		return nil
	}

	var resCount int32
	queue := make(chan *policyQueueItem, len(roles))
	threadNum := a.app.Config().Casbin.LoadThread
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
		list, err := a.menuSVC.List(ctx, dto.MenuListReq{
			RoleID: item.ID,
			Type:   po.MenuType_BUTTON,
			Pager: dto.Pager{
				Page: -1,
			},
		})
		if err != nil {
			logger.Error(ctx, "Failed to query role menus", err)
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
		policyFile := a.app.Config().Casbin.GenPolicyFile
		_ = os.Rename(policyFile, policyFile+".bak")
		_ = os.MkdirAll(filepath.Dir(policyFile), 0755)
		if err := os.WriteFile(policyFile, buf.Bytes(), 0666); err != nil {
			logger.Error(ctx, "Failed to write policy file", err)
			return err
		}
		// set readonly
		_ = os.Chmod(policyFile, 0444)

		modelFile := a.app.Config().Casbin.ModelFile
		e, err := casbin.NewEnforcer(modelFile, policyFile)
		if err != nil {
			logger.Error(ctx, "Failed to create casbin enforcer", err)
			return err
		}
		e.EnableLog(a.app.Config().IsDebug())
		a.enforcer.Store(e)
	}

	logger.Info(ctx, "Casbin load policy",
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
	a.ticker = time.NewTicker(time.Duration(a.app.Config().Casbin.AutoLoadInterval) * time.Second)
	for range a.ticker.C {
		updated, err := a.roleSVC.GetUpdateTime(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to get role update time", err)

			if err := a.roleSVC.RefreshUpdateTime(ctx); err != nil {
				panic(err)
			}
			continue
		}

		if lastUpdated < updated {
			if err := a.load(ctx); err != nil {
				logger.Error(ctx, "Failed to load casbin policy", err)
			} else {
				lastUpdated = updated
			}
		}
	}
}

func (a *Casbinx) Release() error {
	if a.ticker != nil {
		a.ticker.Stop()
	}
	return nil
}
