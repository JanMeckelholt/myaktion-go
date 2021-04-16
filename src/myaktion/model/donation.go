package model

type Status string

const (
	StatusInProcess   Status = "IN_PROCESS"
	StatusTransferred Status = "TRANSFERRED"
)

type Donation struct {
	Amount           float64 `json:"amount"`
	ReceiptRequested bool    `json:"receiptRequested"`
	DonorName        string  `json:"donorName"`
	Status           Status  `json:"status"`
	Account          Account `json:"account"`
}
