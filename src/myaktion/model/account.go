package model

type Account struct {
	Iban       string `json:"iban" gorm:"notNull;size:20"`
	Name       string `json:"name" gorm:"notNull;size:60"`
	NameOfBank string `json:"nameOfBank" gorm:"notNull;size:40"`
}
