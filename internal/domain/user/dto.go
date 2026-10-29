package user

type CreateUserRequest struct {
	Name string `json:"name" form:"name" validate:"required,min=2,max=100"`
}

type UpdateUserRequest struct {
	Name   string `json:"name" form:"name" validate:"required,min=2,max=100"`
	Status string `json:"status" form:"status" validate:"omitempty,oneof=active inactive"`
}
