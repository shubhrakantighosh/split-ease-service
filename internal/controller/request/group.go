package request

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}
