package dto

type UserRoleListReq struct {
	JoinRole bool     `form:"-"`
	UserIDs  []string `form:"userIds"`
}
