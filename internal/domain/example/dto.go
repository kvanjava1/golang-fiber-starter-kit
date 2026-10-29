package example

type CreateCrudRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" form:"description" validate:"max=500"`
}

type UpdateCrudRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=3,max=100"`
	Description string `json:"description" form:"description" validate:"max=500"`
	Status      string `json:"status" form:"status" validate:"omitempty,oneof=active archived inactive"`
}

type SendEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=2"`
}
