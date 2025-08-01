package test

import (
	"gin-admin/internal/dtos"
	"gin-admin/internal/models"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMenu(t *testing.T) {
	e := ApiTester(t)

	t.Cleanup(func() {
		os.RemoveAll("data")
	})

	menuFormItem := dtos.MenuCreateReq{
		Name:  "menu",
		Type:  "menu",
		Path:  "/system/menu",
		Rank:  9,
		Title: "Menu management",
		Extra: map[string]any{
			"icon": "menu",
		},

		Status: models.MenuStatus_ENABLED,
	}

	var menu models.Menu
	e.POST(baseAPI + "/menus").WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&menu))

	assert := assert.New(t)
	assert.NotEmpty(menu.ID)
	assert.Equal(menuFormItem.Name, menu.Name)
	assert.Equal(menuFormItem.Rank, menu.Rank)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Extra, menu.Extra)
	assert.Equal(menuFormItem.Status, menu.Status)

	var getList dtos.ResultList[*models.Menu]
	e.GET(baseAPI + "/menus").Expect().Status(http.StatusOK).JSON().Decode(&getList)
	assert.GreaterOrEqual(len(getList.Data.Items), 1)

	newName := "Menu management 1"
	newStatus := models.MenuStatus_DISABLED
	menu.Name = newName
	menu.Status = newStatus
	e.PUT(baseAPI + "/menus/" + menu.ID).WithJSON(menu).Expect().Status(http.StatusOK)

	var getMenu dtos.Result[*models.Menu]
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK).JSON().Decode(&getMenu)
	assert.Equal(newName, getMenu.Data.Name)
	assert.Equal(newStatus, getMenu.Data.Status)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
