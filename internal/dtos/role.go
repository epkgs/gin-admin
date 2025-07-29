package dtos

// Defining the query parameters for the `Role` struct.
type RoleListReq struct {
	Pager
	Name      string `form:"name"`                                       // Display name of role
	Status    string `form:"status" binding:"oneof=disabled enabled ''"` // Status of role (disabled, enabled
	WithMenus bool   `form:"withMenus"`                                  // Include menu IDs
}

// Defining the data structure for creating a `Role` struct.
type RoleCreateReq struct {
	Code        string   `json:"code" binding:"required,max=32"`                   // Code of role (unique)
	Name        string   `json:"name" binding:"required,max=128"`                  // Display name of role
	Description string   `json:"description"`                                      // Details about role
	Rank        int      `json:"rank"`                                             // Rank for sorting
	Status      string   `json:"status" binding:"required,oneof=disabled enabled"` // Status of role (enabled, disabled)
	MenuIDs     []string `json:"menuIds"`                                          // Menu ids
}

type RoleUpdateReq struct {
	Code        *string   `json:"code" binding:"omitempty,max=32"`                   // Code of role (unique)
	Name        *string   `json:"name" binding:"omitempty,max=128"`                  // Display name of role
	Description *string   `json:"description"`                                       // Details about role
	Rank        *int      `json:"rank"`                                              // Rank for sorting
	Status      *string   `json:"status" binding:"omitempty,oneof=disabled enabled"` // Status of role (enabled, disabled)
	MenuIDs     *[]string `json:"menuIds"`                                           // Menu ids
}
