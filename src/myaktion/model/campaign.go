package model

type Campaign struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	DonationMinimum float64    `json:"donationMinimum"`
	TargetAmount    float64    `json:"targetAmount"`
	Account         Account    `json:"account"`
	OrganizerName   string     `json:"organizerName"`
	Donations       []Donation `json:"donations"`
}
