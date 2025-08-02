package test

import (
	"net/http"
	"os"
	"testing"

	"gin-admin/internal/configs"
	"gin-admin/internal/dtos"
	"gin-admin/internal/models"
	"gin-admin/pkg/crypto/hash"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
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
		Name:  "user",
		Type:  "menu",
		Path:  "/system/user",
		Rank:  7,
		Title: "User management",
		Extra: map[string]any{
			"icon": "user",
		},

		Status: models.MenuStatus_ENABLED,
	}

	var createMenu dtos.Result[*models.Menu]
	e.POST(baseAPI+"/menus").WithHeader("Authorization", "Bearer "+token).WithJSON(menuFormItem).
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

	roleFormItem := dtos.RoleCreateReq{
		Code:        "user",
		Name:        "Normal",
		MenuIDs:     []string{menu.ID},
		Description: "Normal",
		Rank:        8,
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

	var createUser dtos.Result[*models.User]
	e.POST(baseAPI+"/users").WithHeader("Authorization", "Bearer "+token).WithJSON(userFormItem).Expect().Status(http.StatusOK).JSON().Decode(&createUser)
	user := createUser.Data
	assert.NotEmpty(user.ID)
	assert.Equal(userFormItem.Username, user.Username)
	assert.Equal(userFormItem.NickName, user.NickName)
	assert.Equal(userFormItem.Phone, user.Phone)
	assert.Equal(userFormItem.Email, user.Email)
	assert.Equal(userFormItem.Description, user.Description)
	assert.Equal(userFormItem.Status, user.Status)
	assert.Equal(len(userFormItem.RoleIDs), len(user.Roles))

	var listUsers dtos.ResultList[*models.User]
	e.GET(baseAPI+"/users").WithHeader("Authorization", "Bearer "+token).WithQuery("username", userFormItem.Username).Expect().Status(http.StatusOK).JSON().Decode(&listUsers)
	users := listUsers.Data.Items
	assert.GreaterOrEqual(len(users), 1)

	newName := "Test 1"
	newStatus := models.UserStatus_Freezed
	user.NickName = newName
	user.Status = newStatus
	e.PUT(baseAPI+"/users/"+user.ID).WithHeader("Authorization", "Bearer "+token).WithJSON(user).Expect().Status(http.StatusOK)

	var getUser dtos.Result[*models.User]
	e.GET(baseAPI+"/users/"+user.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK).JSON().Decode(&getUser)
	assert.Equal(newName, getUser.Data.NickName)
	assert.Equal(newStatus, getUser.Data.Status)

	e.DELETE(baseAPI+"/users/"+user.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)
	e.GET(baseAPI+"/users/"+user.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI+"/roles/"+role.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)
	e.GET(baseAPI+"/roles/"+role.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)

	e.DELETE(baseAPI+"/menus/"+menu.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusOK)
	e.GET(baseAPI+"/menus/"+menu.ID).WithHeader("Authorization", "Bearer "+token).Expect().Status(http.StatusNotFound)
}
