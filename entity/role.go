package entity

type Role struct {
	SnowflakeID string `json:"snowflake_id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type AddRoleRequest struct {
	SnowflakeId    string   `json:"snowflake_id"`
	Name           string   `json:"name"`
	PermissionList []string `json:"permission_list"`
}

type AddRolePermissionRequest struct {
	SnowflakeId  string `json:"snowflake_id"`
	RoleId       string `json:"role_id"`
	PermissionId string `json:"permission_id"`
}

type UpdateRoleRequest struct {
	SnowflakeId    string   `json:"snowflake_id"`
	Name           string   `json:"name"`
	PermissionList []string `json:"permission_list"`
}
