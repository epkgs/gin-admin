package dtos

// Defining the query parameters for the `RoleMenu` struct.
type MenuRoleListReq struct {
	Pager
	RoleIDs []string `form:"-"` // From Role.ID
	MenuIDs []string `form:"-"` // From Menu.ID
}

// Defining the data structure for creating a `RoleMenu` struct.
type RoleMenuForm struct {
}
