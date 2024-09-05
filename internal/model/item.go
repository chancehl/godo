package model

type GodoItem struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Notes       string `json:"notes"`
	Status      string `json:"status"`
	CompletedOn string `json:"completed_on"`
	CreatedOn   string `json:"created_on"`
}
