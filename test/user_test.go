package test

import (
	"net/http"
	"testing"

	"gin-admin/internal/dtos"
	"gin-admin/internal/models"
	"gin-admin/pkg/crypto/hash"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	e := ApiTester(t)

	menuFormItem := dtos.MenuCreateReq{
		Name: "user",
		Type: "page",
		Path: "/system/user",
		Meta: models.MenuMeta{
			Rank: 7,
			Properties: map[string]any{
				"icon":  "user",
				"title": "User management",
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
		Code:        "user",
		Name:        "Normal",
		MenuIDs:     []string{menu.ID},
		Description: "Normal",
		Rank:        8,
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

	userFormItem := dtos.UserCreateReq{
		Username:    "test",
		NickName:    "Test",
		Password:    hash.MD5String("test"),
		Phone:       "0720",
		Email:       "test@gmail.com",
		Description: "test user",
		Status:      models.UserStatus_Activated,
		RoleIDs:     []string{role.ID},
	}

	var user models.User
	e.POST(baseAPI + "/users").WithJSON(userFormItem).Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&user))
	assert.NotEmpty(user.ID)
	assert.Equal(userFormItem.Username, user.Username)
	assert.Equal(userFormItem.NickName, user.NickName)
	assert.Equal(userFormItem.Phone, user.Phone)
	assert.Equal(userFormItem.Email, user.Email)
	assert.Equal(userFormItem.Description, user.Description)
	assert.Equal(userFormItem.Status, user.Status)
	assert.Equal(len(userFormItem.RoleIDs), len(user.Roles))

	var users models.Users
	e.GET(baseAPI+"/users").WithQuery("username", userFormItem.Username).Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&users))
	assert.GreaterOrEqual(len(users), 1)

	newName := "Test 1"
	newStatus := models.UserStatus_Freezed
	user.NickName = newName
	user.Status = newStatus
	e.PUT(baseAPI + "/users/" + user.ID).WithJSON(user).Expect().Status(http.StatusOK)

	var getUser models.User
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK).JSON().Decode(dtos.NewResultData(&getUser))
	assert.Equal(newName, getUser.NickName)
	assert.Equal(newStatus, getUser.Status)

	e.DELETE(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/users/" + user.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/roles/" + role.ID).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusOK)
	e.GET(baseAPI + "/menus/" + menu.ID).Expect().Status(http.StatusNotFound)
}
