package service

import (
	log "github.com/sirupsen/logrus"

	"github.com/JanMeckelholt/myaktion-go/src/myaktion/model"
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
