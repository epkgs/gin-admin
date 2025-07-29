package test

import (
	"net/http"
	"testing"

	"gin-admin/internal/dtos"
	"gin-admin/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	e := ApiTester(t)

	menuFormItem := dtos.MenuCreateReq{
		Name: "menu",
		Type: "page",
		Path: "/system/menu",
		Meta: models.MenuMeta{
			Rank: 9,
			Properties: map[string]any{
				"icon":  "menu",
				"title": "Menu management",
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

	var menus models.Menus
	e.GET(baseAPI + "/menus").Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&menus))
	assert.GreaterOrEqual(len(menus), 1)

	newName := "Menu management 1"
	newStatus := models.MenuStatus_DISABLED
	menu.Name = newName
	menu.Status = newStatus
	e.PUT(baseAPI + "/menus/" + menu.ID).WithJSON(menu).Expect().Status(http.StatusOK)

	var getMenu models.Menu
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&getMenu))
	assert.Equal(newName, getMenu.Name)
	assert.Equal(newStatus, getMenu.Status)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
