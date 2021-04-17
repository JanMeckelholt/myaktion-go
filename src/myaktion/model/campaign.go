package model

import "gorm.io/gorm"

type Campaign struct {
	gorm.Model
	//ID              uint       `gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"notNull;size:30"`
	DonationMinimum float64    `json:"donationMinimum" gorm:"notNull;check:donation_minimum>=1.0"`
	TargetAmount    float64    `json:"targetAmount" gorm:"notNull;check:target_amount >= 10.0"`
	Account         Account    `gorm:"embedded;embeddedPrefix:account_"`
	OrganizerName   string     `json:"organizerName" gorm:"notNull"`
	Donations       []Donation `json:"donations" gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
