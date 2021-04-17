package handler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/JanMeckelholt/myaktion-go/src/myaktion/model"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/service"
)

type result struct {
	Success string `json:"success"`
}

func CreateCampaign(w http.ResponseWriter, r *http.Request) {
	var campaign *model.Campaign
	campaign, err := requestToCampaign(r)

	if err != nil {
		log.Errorf("Can't serialize request body to campaign struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.CreateCampaign(campaign); err != nil {
		log.Errorf("Error calling service CreateCampaign: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, campaign)
}

func GetCampaigns(w http.ResponseWriter, _ *http.Request) {
	campaigns, err := service.GetCampaigns()
	if err != nil {
		log.Errorf("Error calling serivce GetCampaigns: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, campaigns)
}

func GetCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		log.Errorf("Error getting Id: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	campaign, err := service.GetCampaignById(id)
	if err != nil {
		log.Errorf("Error calling serivce GetCampaignById: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, *campaign)

}

func UpdateCampaign(w http.ResponseWriter, r *http.Request) {
	var newCampaign *model.Campaign
	newCampaign, err := requestToCampaign(r)
	if err != nil {
		log.Errorf("Can't serialize request body to campaign struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := getId(r)
	if err != nil {
		log.Errorf("Error getting Id: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newCampaign, err = service.UpdateCampaignById(id, newCampaign)
	if err != nil {
		log.Errorf("Error calling serivce UpdateCampaignById: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, *newCampaign)
}

func DeleteCampaign(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		log.Errorf("Error getting Id: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	deletedCampaign, err := service.DeleteCampaignById(id)
	if err != nil {
		log.Errorf("Error calling serivce DeleteCampaignById: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, *deletedCampaign)
}

func requestToCampaign(r *http.Request) (*model.Campaign, error) {
	var campaign model.Campaign
	err := json.NewDecoder(r.Body).Decode(&campaign)
	if err != nil {
		log.Errorf("Can't serialize request body to campaign struct: %v", err)
		return nil, err
	}
	return &campaign, nil
}
