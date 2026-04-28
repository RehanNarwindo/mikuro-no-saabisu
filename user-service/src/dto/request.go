package dto

type GetUserByIdRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid"`
}

type GetAllUserHandlersRequest struct {
	Search  string `json:"search" form:"search"`     
	Role    string `json:"role" form:"role"`        
	Limit   int    `json:"limit" form:"limit"`       
	Offset  int    `json:"offset" form:"offset"`     
	SortBy  string `json:"sort_by" form:"sort_by"`   
	SortDir string `json:"sort_dir" form:"sort_dir"` 
}


type CreateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdateUserRequest struct {
	Email     string `json:"email" binding:"omitempty,email"`
	FirstName string `json:"first_name" binding:"omitempty"`
	LastName  string `json:"last_name" binding:"omitempty"`
}