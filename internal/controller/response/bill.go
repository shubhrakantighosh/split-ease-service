package response

type Bill struct {
	User        User    `json:"user"`
	PaidAmount  float64 `json:"paid_amount"`
	Description string  `json:"description"`
}

type Bills []Bill
