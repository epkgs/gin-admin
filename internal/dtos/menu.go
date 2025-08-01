package dtos

// Defining the query parameters for the `Menu` struct.
type MenuListReq struct {
	Pager
	LikeName         string   `form:"name"`          // Display name of menu
	InIDs            []string `form:"-"`             // Include menu IDs
	Status           string   `form:"status"`        // Status of menu (disabled, enabled)
	ParentID         string   `form:"-"`             // Parent ID (From Menu.ID)
	ParentPathPrefix string   `form:"-"`             // Parent path (split by .)
	UserID           string   `form:"-"`             // User ID
	RoleID           string   `form:"-"`             // Role ID
	WithResources    bool     `form:"withResources"` // Include resources
	Type             string   `form:"type"`          // Type of menu (catalog, menu, button)
}

// Defining the data structure for creating a `Menu` struct.
type MenuCreateReq struct {
	Name      string `json:"name" binding:"required_unless=Type button,max=128"` // Display name of menu
	Type      string `json:"type" binding:"required,oneof=catalog menu button"`  // Type of menu (catalog menu, button)
	Path      string `json:"path"`                                               // Access path of menu
	Component string `json:"component"`                                          // Component path of view
	Status    string `json:"status" binding:"required,oneof=disabled enabled"`   // Status of menu (enabled, disabled)
	ParentID  string `json:"parentId"`                                           // Parent ID (From Menu.ID)
	Method    string `json:"method"`                                             // Http method of resource

	Meta map[string]any `json:"meta"` // Meta of menu (JSON)
}

type MenuUpdateReq struct {
	Name      *string `json:"name" binding:"omitempty,max=128"`                   // Display name of menu
	Type      *string `json:"type" binding:"omitempty,oneof=catalog menu button"` // Type of menu (catalog menu, button)
	Path      *string `json:"path"`                                               // Access path of menu
	Component *string `json:"component"`                                          // Component path of view
	Status    *string `json:"status" binding:"omitempty,oneof=disabled enabled"`  // Status of menu (enabled, disabled)
	ParentID  *string `json:"parentId"`                                           // Parent ID (From Menu.ID)
	Method    *string `json:"method"`                                             // Http method of resource

	Meta map[string]any `json:"meta"` // Meta of menu (JSON)
}
