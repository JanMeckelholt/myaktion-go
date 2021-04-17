package service

import (
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/db"
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
	/*campaign.ID = actCampaignId
	campaignStore[actCampaignId] = campaign
	actCampaignId += */
	result := db.DB.Create(campaign)
	if result.Error != nil {
		return result.Error
	}
	log.Infoln("Successfully stored new campaign with ID %v in database.", campaign.ID)
	log.Infoln("Stored: %v", campaign)
	return nil
}

func GetCampaigns() ([]model.Campaign, error) {
	var campaigns []model.Campaign
	result := db.DB.Preload("Donations").Find(&campaigns)
	if result.Error != nil {
		return nil, result.Error
	}
	/*	for _, campaign := range campaignStore {
		campaigns = append(campaigns, *campaign)
	}*/
	log.Tracef("Retrieved: %v", campaigns)
	return campaigns, nil
}

func GetCampaignById(id uint) (*model.Campaign, error) {
	campaign := new(model.Campaign)
	//var campaign *model.Campaign
	log.Infoln("Entered GetCampaignById")
	log.Infof("id: %v", id)
	result := db.DB.Preload("Donations").First(campaign, id)
	//result := db.DB.Take(campaign)
	if result.Error != nil {
		log.Infoln(result)
		return nil, result.Error
	}
	return campaign, nil
}

func UpdateCampaignById(id uint, campaign *model.Campaign) (*model.Campaign, error) {
	existingCampaign := new(model.Campaign)
	result := db.DB.First(existingCampaign, id)
	if result.Error != nil {
		return nil, result.Error
	}
	existingCampaign.Name = campaign.Name
	existingCampaign.OrganizerName = campaign.OrganizerName
	existingCampaign.TargetAmount = campaign.TargetAmount
	existingCampaign.DonationMinimum = campaign.DonationMinimum
	existingCampaign.Account = campaign.Account
	result = db.DB.Save(existingCampaign)
	if result.Error != nil {
		return nil, result.Error
	}
	entry := log.WithField("ID", id)
	entry.Info("Successfully updated campaign.")
	entry.Tracef("Updated: %v", existingCampaign)
	return existingCampaign, nil

}

func DeleteCampaignById(id uint) (*model.Campaign, error) {
	campaign := new(model.Campaign)
	result := db.DB.Delete(campaign, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return campaign, nil
}
