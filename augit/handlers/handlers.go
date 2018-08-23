package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/saml/samlsp"
)

type ShowAccountsResponse struct {
	Users           []*models.GithubUser     `json:"users"`
	ServiceAccounts []*models.ServiceAccount `json:"service_accounts"`
}

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

func ShowAccounts(ghudb models.GithubUserAccessor, sadb models.ServiceAccountAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users, err := ghudb.ListGHUsers()
		if err != nil {
			log.Printf("Failed to find users. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find users"}`))
			return
		}
		serviceAccounts, err := sadb.List()
		if err != nil {
			log.Printf("Failed to find service accounts. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find service accounts"}`))
			return
		}
		allAccounts := ShowAccountsResponse{
			Users:           users,
			ServiceAccounts: serviceAccounts,
		}
		marshaledAccounts, err := json.Marshal(allAccounts)
		if err != nil {
			log.Printf("Failed to marshall account data. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not marshal account data"}`))
			return
		}
		w.Write(marshaledAccounts)
	}
}

type addGHRequest struct {
	GithubID string `json:"github_id"`
}

func AddUser(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req addGHRequest
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
		err = ghudb.ReplaceGHRow(updateUser)
		if err != nil {
			log.Printf("Failed to create user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not create user"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func CheckAdmin(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := getEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(email)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			log.Printf("Failed to find user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}
		if !user.Admin {
			w.Header().Set("Content-Type", "application/json")
			log.Printf("Non-admin attempted to add service account: %s", user.Email)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "User is not admin"}`))
			return
		}
		w.WriteHeader(200)
	}
}

func AddServiceAccount(ghudb models.GithubUserAccessor, sadb models.ServiceAccountAccessor) func(w http.ResponseWriter, r *http.Request) {
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
		if !user.Admin {
			log.Printf("Non-admin attempted to add service account: %s", user.Email)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Must be admin to add service account"}`))
			return
		}

		var req addGHRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Failed to parse user data. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Could not parse user payload"}`))
			return
		}

		newSA := &models.ServiceAccount{
			GithubID:         req.GithubID,
			AdminResponsible: user.ID,
		}
		err = sadb.Upsert(newSA)
		if err != nil {
			log.Printf("Failed to create service account. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not create service account"}`))
			return
		}
		err = ghudb.Delete(req.GithubID)
		if err != nil {
			log.Printf("Failed to delete existing GitHub account record. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not delete existing GitHub account record"}`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
