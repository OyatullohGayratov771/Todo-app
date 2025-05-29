package models

type TaskResponse struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    Done        bool   `json:"done"`
    UserID      string `json:"user_id"`
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}