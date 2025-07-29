package dtos

// Defining the query parameters for the `User` struct.
type UserListReq struct {
	Pager
	LikeUsername string `form:"username"`                                           // Username for login
	LikeName     string `form:"name"`                                               // Name of user
	Status       string `form:"status" binding:"omitempty,oneof=activated freezed"` // Status of user (activated, freezed)
	WithRoles    bool   `form:"withRoles"`                                          // Whether to include role IDs
}

// Defining the data structure for creating a `User` struct.
type UserCreateReq struct {
	Username    string   `json:"username" binding:"required,max=64"`                // Username for login
	NickName    string   `json:"nickName" binding:"required,max=64"`                // Name of user
	RealName    string   `json:"realName" binding:"max=64"`                         // Real name of user
	Password    string   `json:"password" binding:"max=64"`                         // Password for login (md5 hash)
	Wechat      string   `json:"wechat" binding:"max=64"`                           // Wechat account
	Phone       string   `json:"phone" binding:"max=32"`                            // Phone number of user
	Email       string   `json:"email" binding:"omitempty,max=128,email"`           // Email of user
	Description string   `json:"description" binding:"max=1024"`                    // Description of user
	Status      string   `json:"status" binding:"required,oneof=activated freezed"` // Status of user (activated, freezed)
	RoleIDs     []string `json:"roles" binding:"required"`                          // Roles of user
}

type UserUpdateReq struct {
	Username    *string   `json:"username" binding:"omitempty,max=64"`                // Username for login
	NickName    *string   `json:"nickName" binding:"omitempty,max=64"`                // Name of user
	RealName    *string   `json:"realName" binding:"omitempty,max=64"`                // Real name of user
	Password    *string   `json:"password" binding:"omitempty,max=64"`                // Password for login (md5 hash)
	Wechat      *string   `json:"wechat" binding:"omitempty,max=64"`                  // Wechat account
	Phone       *string   `json:"phone" binding:"omitempty,max=32"`                   // Phone number of user
	Email       *string   `json:"email" binding:"omitempty,email,max=128"`            // Email of user
	Description *string   `json:"description" binding:"omitempty,max=1024"`           // Description of user
	Status      *string   `json:"status" binding:"omitempty,oneof=activated freezed"` // Status of user (activated, freezed)
	RoleIDs     *[]string `json:"roles" binding:"omitempty"`                          // Roles of user
}
