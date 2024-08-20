package model

type GodoItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Notes       string `json:"notes"`
	Status      string `json:"status"`
	CompletedOn string `json:"completed_on"`
	CreatedOn   string `json:"created_on"`
}
