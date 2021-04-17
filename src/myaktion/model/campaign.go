package model

import "gorm.io/gorm"

type Campaign struct {
	gorm.Model
	//ID              uint       `gorm:"primaryKey"`
	Name               string     `json:"name" gorm:"notNull;size:30"`
	DonationMinimum    float64    `json:"donationMinimum" gorm:"notNull;check:donation_minimum>=1.0"`
	TargetAmount       float64    `json:"targetAmount" gorm:"notNull;check:target_amount >= 10.0"`
	Account            Account    `gorm:"embedded;embeddedPrefix:account_"`
	OrganizerName      string     `json:"organizerName" gorm:"notNull"`
	Donations          []Donation `json:"donations" gorm:"foreignKey:CampaignID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AmountDonatedSoFar float64    `gorm:"-"`
}

func (c *Campaign) AfterFind(tx *gorm.DB) (err error) {
	var sum float64
	result := tx.Model(&Donation{}).Select("ifnull(sum(amount),0)").Where("campaign_id = ?", c.ID).Scan(&sum)
	if result.Error != nil {
		return result.Error
	}
	c.AmountDonatedSoFar = sum
	return nil
}
