package response

type Bill struct {
	User        User             `json:"user"`
	PaidAmount  float64          `json:"paid_amount"`
	Description string           `json:"description"`
	Splits      []BillSplitEntry `json:"splits,omitempty"` // Who owes this bill payer
}

type Bills []Bill

type BillSplitEntry struct {
	FromUser  User    `json:"from_user"`  // Who owes
	AmountDue float64 `json:"amount_due"` // How much
	IsPaid    bool    `json:"is_paid"`    // If settled
}
