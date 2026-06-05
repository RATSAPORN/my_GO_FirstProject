package dtos

type UserDTO struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Role      string `json:"role"`
}

type NewUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email"    binding:"required,email"`
}

type GetAllUsersRequest struct {
	Page    int `form:"page"`
	PerPage int `form:"per_page"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Email    string `json:"email"    binding:"required,email"`
	Role     string `json:"role" binding:"omitempty,oneof=admin user"`
}

type CommonPaginationData struct {
	Data   any `json:"data"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
