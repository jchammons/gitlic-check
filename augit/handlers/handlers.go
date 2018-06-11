package handlers

import (
	"encoding/json"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/saml/samlsp"
	"log"
	"net/http"
)

func getEmail(token *samlsp.AuthorizationToken) string {
	return token.Subject
}

func ShowUser(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		email := getEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(email)
		if err != nil {
			log.Printf("Failed to find user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}
		marshaledUser, err := json.Marshal(user)
		if err != nil {
			log.Printf("Failed to marshall user data. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}
		w.Write(marshaledUser)
	}
}

type addUserRequest struct {
	GithubID string `json:"github_id"`
}

func AddUser(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req addUserRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Failed to parse user data. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Could not parse user payload"}`))
			return
		}

		updateUser := &models.GithubUser{
			Email:    getEmail(samlsp.Token(r.Context())),
			GithubID: req.GithubID,
		}
		err = ghudb.Upsert(updateUser)
		if err != nil {
			log.Printf("Failed to create user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error:" "Could not create user"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
