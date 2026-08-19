package models

type Draft struct {
	ID          int    `json:"id"`
	EntityName  string `json:"entity_name"`
	EntityID    int    `json:"entity_id"`
	JSONChanges string `json:"json_changes"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
