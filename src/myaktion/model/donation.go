package model

import "gorm.io/gorm"

type Status string

const (
	StatusInProcess   Status = "IN_PROCESS"
	StatusTransferred Status = "TRANSFERRED"
)

type Donation struct {
	gorm.Model
	CampaignID       uint
	Amount           float64 `json:"amount" gorm:"notNull;check:amount>=1.0"`
	ReceiptRequested bool    `json:"receiptRequested"`
	DonorName        string  `json:"donorName" gorm:"notNull;size:40"`
	Status           Status  `json:"status" gorm:"notNull,type:ENUM('TRANSFERED', 'IN_PROCESS')"`
	Account          Account `gorm:"embedded;embeddedPrefix:account_"`
}
