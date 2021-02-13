package services

type EmailData struct {
	To   string `json:"to" form:"to" binding:"required,email"`
	Body string `json:"body" form:"body" binding:"required"`
}
