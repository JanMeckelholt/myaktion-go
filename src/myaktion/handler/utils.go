package handler

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func getId(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		log.Errorf("Can't get ID from request: %v", err)
		return 0, err
	}
	return uint(id), nil
}
