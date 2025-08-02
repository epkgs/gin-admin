package test

import (
	"net/http"
	"os"
	"testing"

	"gin-admin/internal/configs"
	"gin-admin/internal/dtos"
	"gin-admin/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	e := ApiTester(t)

	t.Cleanup(func() {
		os.RemoveAll("data")
	})

	var login dtos.Result[*dtos.LoginToken]
	e.POST(baseAPI + "/auth/login").WithJSON(dtos.Login{
		Username: configs.C.Super.Username,
		Password: configs.C.Super.Password,
	}).Expect().Status(http.StatusOK).JSON().Decode(&login)

	token := login.Data.AccessToken

	menuFormItem := dtos.MenuCreateReq{
		Name:  "role",
		Type:  "menu",
		Path:  "/system/role",
		Rank:  8,
		Title: "Role management",
		Extra: map[string]any{
			"icon": "role",
		},

		Status: models.MenuStatus_ENABLED,
	}

	var createMenu dtos.Result[*models.Menu]
	e.POST(baseAPI+"/menus").WithHeader("Authorization", "Bearer "+token).WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(&createMenu)

	menu := createMenu.Data

	assert := assert.New(t)
	assert.NotEmpty(menu.ID)
	assert.Equal(menuFormItem.Name, menu.Name)
	assert.Equal(menuFormItem.Rank, menu.Rank)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Extra, menu.Extra)
	assert.Equal(menuFormItem.Status, menu.Status)

	roleFormItem := dtos.RoleCreateReq{
		Code:        "admin",
		Name:        "Administrator",
		MenuIDs:     []string{menu.ID},
		Description: "Administrator",
		Rank:        9,
		Status:      models.RoleStatus_Enabled,
	}

	var createRole dtos.Result[*models.Role]
	e.POST(baseAPI+"/roles").WithHeader("Authorization", "Bearer "+token).WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(&createRole)
	role := createRole.Data
	assert.NotEmpty(role.ID)
	assert.Equal(roleFormItem.Code, role.Code)
	assert.Equal(roleFormItem.Name, role.Name)
	assert.Equal(roleFormItem.Description, role.Description)
	assert.Equal(roleFormItem.Rank, role.Rank)
	assert.Equal(roleFormItem.Status, role.Status)
	assert.Equal(len(roleFormItem.MenuIDs), len(role.Menus))

	var listRoles dtos.ResultList[*models.Role]
	e.GET(baseAPI+"/roles").WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Decode(&listRoles)
	roles := listRoles.Data.Items
	assert.GreaterOrEqual(len(roles), 1)

	newName := "Administrator 1"
	newStatus := models.RoleStatus_Disabled
	role.Name = newName
	role.Status = newStatus
	e.PUT(baseAPI+"/roles/"+role.ID).WithHeader("Authorization", "Bearer "+token).WithJSON(role).Expect().Status(http.StatusOK)

	var getRole dtos.Result[*models.Role]
	e.GET(baseAPI+"/roles/"+role.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Decode(&getRole)
	assert.Equal(newName, getRole.Data.Name)
	assert.Equal(newStatus, getRole.Data.Status)

	e.DELETE(baseAPI+"/roles/"+role.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)
	e.GET(baseAPI+"/roles/"+role.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI+"/menus/"+menu.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)
	e.GET(baseAPI+"/menus/"+menu.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)
}
