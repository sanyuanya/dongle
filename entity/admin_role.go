package entity

type AddAdminRoleRequest struct {
	SnowflakeId string `json:"snowflake_id"`
	AdminId     string `json:"admin_id"`
	RoleId      string `json:"role_id"`
}

type GetAdminRoleResponse struct {
	AdminId   string `json:"admin_id"`
	RoleId    string `json:"role_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
