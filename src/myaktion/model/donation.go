package model

type Status string

const (
	StatusInProcess   Status = "IN_PROCESS"
	StatusTransferred Status = "TRANSFERRED"
)

type Donation struct {
	Amount           float64
	ReceiptRequested bool
	DonorName        string
	Status           Status
	Account          Account
}
