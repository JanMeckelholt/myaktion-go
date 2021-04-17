package handler

import (
	"encoding/json"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/model"
	"github.com/JanMeckelholt/myaktion-go/src/myaktion/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AddDonation(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		log.Errorf("Error getting Id: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	donation, err := requestToDonation(r)
	if err != nil {
		log.Errorf("Can't serialize request body to donation struct: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = service.AddDonation(id, donation)
	if err != nil {
		log.Errorf("Error calling serivce AddDonation: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//sendJson(w, *updatedCampaign)
}

func requestToDonation(r *http.Request) (*model.Donation, error) {
	var donation model.Donation
	err := json.NewDecoder(r.Body).Decode(&donation)
	if err != nil {
		log.Errorf("Can't serialize request body to donation struct: %v", err)
		return nil, err
	}
	return &donation, nil
}
