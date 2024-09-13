package entity

type Permission struct {
	SnowflakeId string `json:"snowflake_id"`
	Summary     string `json:"summary"`
	Path        string `json:"path"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PermissionMenu struct {
	SnowflakeId string            `json:"snowflake_id"`
	Summary     string            `json:"summary"`
	Path        string            `json:"path"`
	ApiPath     string            `json:"api_path"`
	ParentId    string            `json:"parent_id"`
	Children    []*PermissionMenu `json:"children"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}
