package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/solarwinds/gitlic-check/augit"
	"github.com/solarwinds/gitlic-check/augit/email"
	"github.com/solarwinds/gitlic-check/augit/models"
	"github.com/solarwinds/saml/samlsp"
)

const (
	ERR_DB_RETRIEVAL = 50003
	ERR_MARSHAL_DATA = 50004
	ERR_BAD_INPUT    = 50005
	ERR_DB_WRITE     = 50006
	ERR_NOT_ADMIN    = 50007
	ERR_ID_TAKEN     = 50008
	ERR_DB_DELETE    = 50009
	ERR_FORBIDDEN    = 50010
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		marshaledUser, err := json.Marshal(user)
		if err != nil {
			log.Printf("Failed to marshall user data. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_MARSHAL_DATA, "Could not marshal response")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find users")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		serviceAccounts, err := sadb.List()
		if err != nil {
			log.Printf("Failed to find service accounts. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find service accounts")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		svcAcctResponses := []*serviceAcctResponse{}
		for _, svcAcct := range serviceAccounts {
			admin, err := ghudb.FindByID(svcAcct.AdminResponsible)
			if err != nil {
				augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user for service account")
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
			errPayload := augit.LogAndFormatError(ERR_MARSHAL_DATA, "Could not marshal response")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		w.Write(marshaledAccounts)
	}
}

func ShowLog(ldb models.AuditLogAccessor) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type logResponse struct {
			Entries []*models.AuditLog `json:"entries"`
		}
		entries, err := ldb.List()
		if err != nil {
			log.Printf("Failed to find audit log entries. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find audit log entries")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		resp := &logResponse{entries}
		marshaledEntries, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Failed to marshall audit log data. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_MARSHAL_DATA, "Could not marshal response")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaledEntries)
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
			errPayload := augit.LogAndFormatError(ERR_BAD_INPUT, "Could not parse user input")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_WRITE, "Could not create user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		if !user.Admin {
			w.Header().Set("Content-Type", "application/json")
			log.Printf("Non-admin attempted to add service account: %s", user.Email)
			errPayload := augit.LogAndFormatError(ERR_NOT_ADMIN, "User is not admin")
			w.WriteHeader(http.StatusForbidden)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		var req ghRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Printf("Failed to parse user data. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_BAD_INPUT, "Could not parse user input")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
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
				errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not verify GitHub ID's current registration status")
				w.WriteHeader(http.StatusBadGateway)
				w.Write(errPayload)
				return
			}
			if ghEntry.Username != "" || ghEntry.Email != "" {
				errPayload := augit.LogAndFormatError(ERR_ID_TAKEN, fmt.Sprintf("GitHub ID %s is already registered to a SolarWinds user", req.GithubID))
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errPayload)
				return
			}

			err = sadb.Create(newSA)
			if err != nil {
				log.Printf("Failed to create service account. Error: %v\n", err)
				errPayload := augit.LogAndFormatError(ERR_DB_WRITE, "Could not create service account")
				w.WriteHeader(http.StatusBadGateway)
				w.Write(errPayload)
				return
			}
			err = ghudb.Delete(req.GithubID)
			if err != nil && !models.IsErrRecordNotFound(err) {
				log.Printf("Failed to delete existing GitHub account record. Error: %v\n", err)
				errPayload := augit.LogAndFormatError(ERR_DB_DELETE, "Could not delete existing GitHub account record")
				w.WriteHeader(http.StatusBadGateway)
				w.Write(errPayload)
				return
			}

		} else {
			log.Printf("Attempt to add existing service account: %s\n", req.GithubID)
			errPayload := augit.LogAndFormatError(ERR_ID_TAKEN, "Service account already registered")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
			return
		}
		owners, err := ghodb.List()
		if err != nil {
			log.Printf("Failed to find owners. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find owners")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}

		githubID, ok := vars["githubid"]
		if !ok {
			errPayload := augit.LogAndFormatError(ERR_BAD_INPUT, "Must supply the GitHub ID for the service account you want to remove")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
			return
		}

		sa, err := sadb.FindByGithubID(githubID)
		if err != nil {
			log.Printf("Failed to find service account for deletion. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find service account")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}

		if sa.AdminResponsible != user.ID {
			isOwner, err := ghodb.ExistsByGithubID(user.GithubID)
			if err != nil {
				log.Printf("Failed to verify if submitter is a GitHub owner. Error: %v\n", err)
			}
			if !isOwner {
				log.Printf("Only the user who registered a service account or a GitHub org owner may remove it")
				errPayload := augit.LogAndFormatError(ERR_FORBIDDEN, "Only the user who registered a service account or a GitHub org owner may remove it")
				w.WriteHeader(http.StatusForbidden)
				w.Write(errPayload)
				return
			}
		}

		err = sadb.Delete(githubID)
		if err != nil {
			log.Printf("Failed to remove service account. Error: %v\n", err)
			errPayload := augit.LogAndFormatError(ERR_DB_DELETE, "Could not delete existing service account")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		if !user.Admin {
			log.Printf("Non-admin attempted to add admin: %s", user.Email)
			errPayload := augit.LogAndFormatError(ERR_NOT_ADMIN, "Must be admin to add admin")
			w.WriteHeader(http.StatusForbidden)
			w.Write(errPayload)
			return
		}

		adminEmail, ok := vars["email"]
		if !ok {
			log.Printf("Email not provided in AddAdmin URL")
			errPayload := augit.LogAndFormatError(ERR_BAD_INPUT, "Email not provided in URL")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
			return
		}
		err = ghudb.AddAdmin(adminEmail)
		if err != nil {
			if models.IsErrRecordNotFound(err) {
				log.Printf("error adding admin %s; email not registered\n", adminEmail)
				errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Requested admin is not a valid user")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errPayload)
				return
			}
			log.Printf("error adding admin for %s: %s\n", adminEmail, err.Error())
			errPayload := augit.LogAndFormatError(ERR_DB_WRITE, "Could not add admin")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
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
			errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Could not find user")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		if !user.Admin {
			log.Printf("Non-admin attempted to remove admin: %s", user.Email)
			errPayload := augit.LogAndFormatError(ERR_NOT_ADMIN, "Must be admin to remove admin")
			w.WriteHeader(http.StatusForbidden)
			w.Write(errPayload)
			return
		}

		adminEmail, ok := vars["email"]
		if !ok {
			log.Printf("Email not provided in RemoveAdmin URL")
			errPayload := augit.LogAndFormatError(ERR_BAD_INPUT, "Email not provided in URL")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errPayload)
			return
		}
		err = ghudb.RemoveAdmin(adminEmail)
		if err != nil {
			if models.IsErrRecordNotFound(err) {
				log.Printf("error adding admin %s; email not registered\n", adminEmail)
				errPayload := augit.LogAndFormatError(ERR_DB_RETRIEVAL, "Requested admin to remove is not a valid user")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errPayload)
				return
			}
			log.Printf("error removing admin for %s: %s\n", adminEmail, err.Error())
			errPayload := augit.LogAndFormatError(ERR_DB_WRITE, "Could not remove admin")
			w.WriteHeader(http.StatusBadGateway)
			w.Write(errPayload)
			return
		}
		log.Printf("admin %s removed admin privileges from %s", user.Email, adminEmail)
		w.WriteHeader(http.StatusCreated)
	}
}
