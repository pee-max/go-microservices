package main

import (
	"encoding/json"
	"net/http"
	"ride-sharing/shared/contracts"
	"time"
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 9)
	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if reqBody.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	// log.Println("SUCCESS")

	response := contracts.APIResponse{Data: "ok"}

	writeJSON(w, http.StatusOK, response)

}
