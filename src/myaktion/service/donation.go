package service

import (
	"context"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/client"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/client/banktransfer"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/db"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/model"
	log "github.com/sirupsen/logrus"
	"time"
)

func AddDonation(campaignId uint, donation *model.Donation) error {
	campaign, err := GetCampaignById(campaignId)
	if err != nil {
		return err
	}
	donation.CampaignID = campaignId
	result := db.DB.Create(donation)
	if result.Error != nil {
		return result.Error
	}
	conn, err := client.GetBankTransferConnection()
	if err != nil {
		log.Errorf("error connecting to the banktransfer service: %v", err)
		deleteDonation(donation)
		return err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	banktransferClient := banktransfer.NewBankTransferClient(conn)
	_, err = banktransferClient.TransferMoney(ctx, &banktransfer.Transaction{
		DonationId:  int32(donation.ID),
		Amount:      float32(donation.Amount),
		Reference:   "Donation",
		FromAccount: convertAccount(&donation.Account),
		ToAccount:   convertAccount(&campaign.Account),
	})
	if err != nil {
		log.Errorf("error calling the banktransfer service: %v", err)
		deleteDonation(donation)
		return err
	}
	entry := log.WithField("ID", campaignId)
	entry.Info("Successfully added new donation to campaign in database.")
	entry.Tracef("Stored: %v", donation)
	return nil
}

func convertAccount(account *model.Account) *banktransfer.Account {
	return &banktransfer.Account{
		Name:     account.Name,
		BankName: account.NameOfBank,
		Number:   account.Iban,
	}
}
func deleteDonation(donation *model.Donation) error {
	entry := log.WithField("donationID", donation.ID)
	entry.Info("Trying to delete donation to make state consistent.")
	result := db.DB.Delete(donation)
	if result.Error != nil {
		// Note: configure logger to raise an alarm to compensate inconsistent state
		entry.WithField("alarm", true).Error("")
		return result.Error
	}
	entry.Info("Successfully deleted campaign.")
	return nil
}
