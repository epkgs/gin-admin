package test

import (
	"net/http"
	"testing"

	"gin-admin/internal/dtos"
	"gin-admin/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	e := ApiTester(t)

	menuFormItem := dtos.MenuCreateReq{
		Name: "role",
		Type: "page",
		Path: "/system/role",
		Meta: models.MenuMeta{
			Rank: 8,
			Properties: map[string]any{
				"icon":  "role",
				"title": "Role management",
			},
		},
		Status: models.MenuStatus_ENABLED,
	}

	var menu models.Menu
	e.POST(baseAPI + "/menus").WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&menu))

	assert := assert.New(t)
	assert.NotEmpty(menu.ID)
	assert.Equal(menuFormItem.Name, menu.Name)
	assert.Equal(menuFormItem.Meta.Rank, menu.Meta.Rank)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Meta, menu.Meta)
	assert.Equal(menuFormItem.Status, menu.Status)

	roleFormItem := dtos.RoleCreateReq{
		Code:        "admin",
		Name:        "Administrator",
		MenuIDs:     []string{menu.ID},
		Description: "Administrator",
		Rank:        9,
		Status:      models.RoleStatus_Enabled,
	}

	var role models.Role
	e.POST(baseAPI + "/roles").WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&role))
	assert.NotEmpty(role.ID)
	assert.Equal(roleFormItem.Code, role.Code)
	assert.Equal(roleFormItem.Name, role.Name)
	assert.Equal(roleFormItem.Description, role.Description)
	assert.Equal(roleFormItem.Rank, role.Rank)
	assert.Equal(roleFormItem.Status, role.Status)
	assert.Equal(len(roleFormItem.MenuIDs), len(role.Menus))

	var roles models.Roles
	e.GET(baseAPI + "/roles").Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&roles))
	assert.GreaterOrEqual(len(roles), 1)

	newName := "Administrator 1"
	newStatus := models.RoleStatus_Disabled
	role.Name = newName
	role.Status = newStatus
	e.PUT(baseAPI + "/roles/" + role.ID).WithJSON(role).Expect().Status(http.StatusOK)

	var getRole models.Role
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&getRole))
	assert.Equal(newName, getRole.Name)
	assert.Equal(newStatus, getRole.Status)

	e.DELETE(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
