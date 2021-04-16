package model

type Campaign struct {
	ID              uint
	Name            string
	DonationMinimum float64
	TargetAmount    float64
	Account         Account
	OrganizerName   string
	Donations       []Donation
}
