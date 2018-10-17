package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/solarwinds/gitlic-check/augit/email"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/saml/samlsp"
)

type serviceAcctResponse struct {
	GithubID         string `json:"github_id"`
	AdminResponsible string `json:"admin_responsible"`
}

type ShowAccountsResponse struct {
	Users           []*models.GithubUser   `json:"users"`
	ServiceAccounts []*serviceAcctResponse `json:"service_accounts"`
}

// getCanonicalEmail accesses the UserName attribute on the SAML token, which thanks to
// our Okta integration is sourced from the Active Directory userPrincipalName field.
// This is the canonical email that we use to relate users in our system back to their
// Active Directory entries
func getCanonicalEmail(token *samlsp.AuthorizationToken) string {
	return token.Attributes.Get("UserName")
}

// getEmail accesses the Subject on the SAML token which should be the user's primary email
// in Okta. We want to use this for actually sending emails, but not for associating a user
// in our system with their Active Directory entry.
func getEmail(token *samlsp.AuthorizationToken) string {
	return token.Subject
}

func ShowUser(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		username := getCanonicalEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(username)
		if err != nil {
			log.Printf("Failed to find user with email %s. Error: %v\n", username, err)
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
		svcAcctResponses := []*serviceAcctResponse{}
		for _, svcAcct := range serviceAccounts {
			admin, err := ghudb.FindByID(svcAcct.AdminResponsible)
			if err != nil {
				continue
			}
			svcAcctResponses = append(svcAcctResponses, &serviceAcctResponse{
				GithubID:         svcAcct.GithubID,
				AdminResponsible: admin.Email,
			})
		}
		allAccounts := ShowAccountsResponse{
			Users:           users,
			ServiceAccounts: svcAcctResponses,
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

type ghRequest struct {
	GithubID string `json:"github_id"`
}

func AddUser(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req ghRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Failed to parse user data. Error: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Could not parse user payload"}`))
			return
		}

		updateUser := &models.GithubUser{
			Email:    getEmail(samlsp.Token(r.Context())),
			Username: getCanonicalEmail(samlsp.Token(r.Context())),
			GithubID: req.GithubID,
		}
		err = ghudb.ReplaceGHRow(updateUser)
		if err != nil {
			log.Printf("Failed to create user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf(`{"error": "Could not create user: %s"}`, err.Error())))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func CheckAdmin(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email := getCanonicalEmail(samlsp.Token(r.Context()))
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

func AddServiceAccount(ghudb models.GithubUserAccessor, ghodb models.GithubOwnerAccessor, sadb models.ServiceAccountAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		inEmail := getCanonicalEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(inEmail)
		if err != nil {
			log.Printf("Failed to find user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}
		var req ghRequest
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
		exists, err := sadb.Exists(req.GithubID)
		if !exists {
			// Check to ensure GH id is not already associated with a SW user
			ghEntry, err := ghudb.FindByGithubID(req.GithubID)
			if err != nil {
				log.Printf("Failed to verify GitHub ID for service account is not already registered. Error: %v\n", err)
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(`{"error": "Could not verify GitHub ID's current registration status.`))
				return
			}
			if ghEntry.Username != "" || ghEntry.Email != "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "GitHub ID is already registered to a SolarWinds user.`))
				return
			}

			err = sadb.Create(newSA)
			if err != nil {
				log.Printf("Failed to create service account. Error: %v\n", err)
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(`{"error": "Could not create service account"}`))
				return
			}
			err = ghudb.Delete(req.GithubID)
			if err != nil && !models.IsErrRecordNotFound(err) {
				log.Printf("Failed to delete existing GitHub account record. Error: %v\n", err)
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte(`{"error": "Could not delete existing GitHub account record"}`))
				return
			}

		} else {
			log.Printf("Attempt to add existing service account: %s\n", req.GithubID)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Service account already registered"}`))
			return
		}
		owners, err := ghodb.List()
		if err != nil {
			log.Printf("Failed to retrieve GitHub owner records. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not retrieve GitHub owner records"}`))
			return
		}
		err = email.SendOwnerListEmail(getEmail(samlsp.Token(r.Context())), req.GithubID, owners)
		w.WriteHeader(http.StatusNoContent)
	}
}

func RemoveServiceAccount(ghudb models.GithubUserAccessor, sadb models.ServiceAccountAccessor, ghodb models.GithubOwnerAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		inEmail := getCanonicalEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(inEmail)
		if err != nil {
			log.Printf("Failed to find user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}

		githubID, ok := vars["githubid"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Must supply the GitHub ID for the service account you want to remove"}`))
			return
		}

		sa, err := sadb.FindByGithubID(githubID)
		if err != nil {
			log.Printf("Failed to find service account for deletion. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find a registered service account with that GitHub ID`))
			return
		}

		if sa.AdminResponsible != user.ID {
			isOwner, err := ghodb.ExistsByGithubID(user.GithubID)
			if err != nil {
				log.Printf("Failed to verify if submitter is a GitHub owner. Error: %v\n", err)
			}
			if !isOwner {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(`{"error": "Only the user who registered a service account or a GitHub org owner may remove it"`))
				return
			}
		}

		err = sadb.Delete(githubID)
		if err != nil {
			log.Printf("Failed to remove service account. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Something went wrong. Could not remove the service account.`))
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func AddAdmin(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := getCanonicalEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(username)
		if err != nil {
			log.Printf("Failed to find user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}
		if !user.Admin {
			log.Printf("Non-admin attempted to add admin: %s", user.Email)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Must be admin to add admin"}`))
			return
		}

		adminEmail, ok := vars["email"]
		if !ok {
			http.Error(w, "must supply an email for new admin", http.StatusBadRequest)
			return
		}
		err = ghudb.AddAdmin(adminEmail)
		if err != nil {
			if models.IsErrRecordNotFound(err) {
				log.Printf("error adding admin %s; email not registered\n", adminEmail)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "User with requested email not found"}`))
				return
			}
			log.Printf("error adding admin for %s: %s\n", adminEmail, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("admin %s added admin privileges to %s", user.Email, adminEmail)
		w.WriteHeader(http.StatusCreated)
	}
}

func RemoveAdmin(ghudb models.GithubUserAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := getCanonicalEmail(samlsp.Token(r.Context()))
		user, err := ghudb.Find(email)
		if err != nil {
			log.Printf("Failed to find user. Error: %v\n", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Could not find user"}`))
			return
		}
		if !user.Admin {
			log.Printf("Non-admin attempted to remove admin: %s", user.Email)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"error": "Must be admin to remove admin"}`))
			return
		}

		adminEmail, ok := vars["email"]
		if !ok {
			http.Error(w, "must supply an email for deleted admin", http.StatusBadRequest)
			return
		}
		err = ghudb.RemoveAdmin(adminEmail)
		if err != nil {
			if models.IsErrRecordNotFound(err) {
				log.Printf("error removing admin %s; email not registered\n", adminEmail)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "User with requested email not found"}`))
				return
			}
			log.Printf("error removing admin for %s: %s\n", adminEmail, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("admin %s removed admin privileges from %s", user.Email, adminEmail)
		w.WriteHeader(http.StatusCreated)
	}
}
