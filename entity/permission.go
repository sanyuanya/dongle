package entity

type Permission struct {
	SnowflakeID string `json:"snowflake_id"`
	Summary     string `json:"summary"`
	Path        string `json:"path"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type PermissionMenu struct {
	SnowflakeID string `json:"snowflake_id"`
	Summary     string `json:"summary"`
	Path        string `json:"path"`
	ApiPath     string `json:"api_path"`
}
