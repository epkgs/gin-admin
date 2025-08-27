package test

import (
	"net/http"
	"os"
	"testing"

	"gin-admin/internal/model/dto"
	"gin-admin/internal/model/po"
	"gin-admin/pkg/crypto/hash"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	e := ApiTester(t)

	t.Cleanup(func() {
		os.RemoveAll("data")
	})

	menuFormItem := dto.MenuCreateReq{
		Name:  "user",
		Type:  "menu",
		Path:  "/system/user",
		Rank:  7,
		Title: "User management",
		Extra: map[string]any{
			"icon": "user",
		},

		Status: po.MenuStatus_ENABLED,
	}

	var createMenu dto.Result[*po.Menu]
	e.POST(baseAPI + "/menus").WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(&createMenu)

	assert := assert.New(t)

	menu := createMenu.Data
	assert.NotEmpty(menu.ID)
	assert.Equal(menuFormItem.Name, menu.Name)
	assert.Equal(menuFormItem.Rank, menu.Rank)
	assert.Equal(menuFormItem.Type, menu.Type)
	assert.Equal(menuFormItem.Path, menu.Path)
	assert.Equal(menuFormItem.Extra, menu.Extra)
	assert.Equal(menuFormItem.Status, menu.Status)

	roleFormItem := dto.RoleCreateReq{
		Code:        "user",
		Name:        "Normal",
		MenuIDs:     []string{menu.ID},
		Description: "Normal",
		Rank:        8,
		Status:      po.RoleStatus_Enabled,
	}

	var createRole dto.Result[*po.Role]
	e.POST(baseAPI + "/roles").WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(&createRole)

	role := createRole.Data
	assert.NotEmpty(role.ID)
	assert.Equal(roleFormItem.Code, role.Code)
	assert.Equal(roleFormItem.Name, role.Name)
	assert.Equal(roleFormItem.Description, role.Description)
	assert.Equal(roleFormItem.Rank, role.Rank)
	assert.Equal(roleFormItem.Status, role.Status)
	assert.Equal(len(roleFormItem.MenuIDs), len(role.Menus))

	userFormItem := dto.UserCreateReq{
		Username:    "test",
		NickName:    "Test",
		Password:    hash.MD5String("test"),
		Phone:       "0720",
		Email:       "test@gmail.com",
		Description: "test user",
		Status:      po.UserStatus_Activated,
		RoleIDs:     []string{role.ID},
	}

	var createUser dto.Result[*po.User]
	e.POST(baseAPI + "/users").WithJSON(userFormItem).Expect().Status(http.StatusOK).JSON().Decode(&createUser)
	user := createUser.Data
	assert.NotEmpty(user.ID)
	assert.Equal(userFormItem.Username, user.Username)
	assert.Equal(userFormItem.NickName, user.NickName)
	assert.Equal(userFormItem.Phone, user.Phone)
	assert.Equal(userFormItem.Email, user.Email)
	assert.Equal(userFormItem.Description, user.Description)
	assert.Equal(userFormItem.Status, user.Status)
	assert.Equal(len(userFormItem.RoleIDs), len(user.Roles))

	var listUsers dto.ResultList[*po.User]
	e.GET(baseAPI+"/users").WithQuery("username", userFormItem.Username).Expect().Status(http.StatusOK).JSON().Decode(&listUsers)
	users := listUsers.Data.Items
	assert.GreaterOrEqual(len(users), 1)

	newName := "Test 1"
	newStatus := po.UserStatus_Freezed
	user.NickName = newName
	user.Status = newStatus
	e.PUT(baseAPI + "/users/" + user.ID).WithJSON(user).Expect().Status(http.StatusOK)

	var getUser dto.Result[*po.User]
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK).JSON().Decode(&getUser)
	assert.Equal(newName, getUser.Data.NickName)
	assert.Equal(newStatus, getUser.Data.Status)

	e.DELETE(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
