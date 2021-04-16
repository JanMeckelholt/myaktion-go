package service

import (
	"errors"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/model"
	log "github.com/sirupsen/logrus"
)

var (
	campaignStore map[uint]*model.Campaign
	actCampaignId uint = 1
)

func init() {
	campaignStore = make(map[uint]*model.Campaign)
}

func CreateCampaign(campaign *model.Campaign) error {
	campaign.ID = actCampaignId
	campaignStore[actCampaignId] = campaign
	actCampaignId += 1
	log.Infoln("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Infoln("Stored: %v", campaign)
	return nil
}

func GetCampaigns() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	for _, campaign := range campaignStore {
		campaigns = append(campaigns, *campaign)
	}
	log.Infoln("Retrieved: %v", campaigns)
	return campaigns, nil
}

func GetCampaignById(id uint) (*model.Campaign, error) {
	//var campaign model.Campaign
	if campaign, ok := campaignStore[id]; ok {
		return campaign, nil
	}

	return nil, errors.New("Campaign for id not found: " + string(id))
}

func UpdateCampaignById(id uint, campaign *model.Campaign) (*model.Campaign, error) {
	if existingCampaign, ok := campaignStore[id]; ok {
		//campaignStore[id] = campaign
		existingCampaign.Name = campaign.Name
		existingCampaign.OrganizerName = campaign.OrganizerName
		existingCampaign.TargetAmount = campaign.TargetAmount
		existingCampaign.DonationMinimum = campaign.DonationMinimum
		existingCampaign.Account = campaign.Account
		entry := log.WithField("ID", id)
		entry.Info("Successfully updated campaign.")
		entry.Tracef("Updated: %v", existingCampaign)
		return existingCampaign, nil
	}
	return nil, errors.New("Campaign for id not found: " + string(id))
}

func DeleteCampaignById(id uint) (*model.Campaign, error) {
	if campaign, ok := campaignStore[id]; ok {
		delete(campaignStore, id)
		return campaign, nil
	}
	return nil, errors.New("Campaign for id not found: " + string(id))
}

func AddDonation(id uint, donation *model.Donation) (*model.Campaign, error) {
	if existingCampaign, ok := campaignStore[id]; ok {
		existingCampaign.Donations = append(existingCampaign.Donations, *donation)
		entry := log.WithField("ID", id)
		entry.Info("Successfully updated campaign.")
		entry.Tracef("Updated: %v", existingCampaign)
		return existingCampaign, nil
	}
	return nil, errors.New("Campaign for id not found: " + string(id))
}
