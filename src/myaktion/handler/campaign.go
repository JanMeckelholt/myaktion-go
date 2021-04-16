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
	campaign, err := getCampaign(r)

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
	/*w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(*campaign); err != nil {
		log.Errorf("Failure encoding value to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
	sendJson(w, campaign)
}

func GetCampaigns(w http.ResponseWriter, _ *http.Request) {
	campaigns, err := service.GetCampaigns()
	if err != nil {
		log.Errorf("Error calling serivce GetCampaings: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*	w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(campaigns); err != nil {
			log.Errorf("Failure encodig value to JSON: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}*/
	sendJson(w, campaigns)
}

func getCampaign(r *http.Request) (*model.Campaign, error) {
	var campaign model.Campaign
	err := json.NewDecoder(r.Body).Decode(&campaign)
	if err != nil {
		log.Errorf("Can't serialize request body to campaign sturct: %v", err)
		return nil, err
	}
	return &campaign, nil
}

func sendJson(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Errorf("Failure encoding value to JSON: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
