package request

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateBillRequest struct {
	PaidAmount  float64 `json:"paid_amount"`
	Description string  `json:"description"`
}

type UpdateBillRequest struct {
	PaidAmount  float64 `json:"paid_amount"`
	Description string  `json:"description"`
}
