package populator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
}

type UserDatabase interface {
	Create(*User) error
}

type MSTokenResponse struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
	Resource    string `json:"resource"`
	AccessToken string `json:"access_token"`
	AccessType  string `json:"access_type"`
}

type AzureADResponse struct {
	Value    []ADUser `json:"value"`
	NextLink string   `json:"odata.nextLink"`
}

type ADUser struct {
	Mail      string `json:"mail"`
	GivenName string `json:"givenName"`
	Surname   string `json:"surname"`
}

type Populator struct {
	clientID     string
	clientSecret string
}

func NewPopulator(clientID, clientSecret string) *Populator {
	return &Populator{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (p *Populator) Populate(userDb UserDatabase) error {
	userResps := p.getAllUsers()
	for _, resp := range userResps {
		for _, user := range resp.Value {
			if user.Mail == "" {
				fmt.Printf("skipping user: %+v\n", user)
				continue
			}
			modelUser := &User{
				FirstName: user.GivenName,
				LastName:  user.Surname,
				Email:     user.Mail,
			}
			err := userDb.Create(modelUser)
			if err != nil {
				// TODO: Create error type for array of errors to keep track of failures
				fmt.Printf("error saving: %+v\n", err)
				fmt.Printf("skipping user: %+v\n", user)
			}
		}
	}
	return nil
}

func (p *Populator) getAllUsers() []*AzureADResponse {
	toke := p.getToken()
	more := true
	nextLinkReg := regexp.MustCompile("skiptoken=(.*)")

	allResps := []*AzureADResponse{}
	adResp := &AzureADResponse{}
	nextLink := ""
	iterations := 1
	for more {
		adResp = p.requestUsers(toke, nextLink)
		allResps = append(allResps, adResp)
		matches := nextLinkReg.FindStringSubmatch(adResp.NextLink)
		if len(matches) == 2 {
			nextLink = matches[1]
		} else {
			nextLink = ""
		}
		if adResp.NextLink == "" {
			more = false
		}
		iterations += 1
	}
	return allResps
}

func (p *Populator) requestUsers(token *MSTokenResponse, skipToken string) *AzureADResponse {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.windows.net/solarwinds.com/users?api-version=1.6&$top=999", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Referer", "https://graphexplorer.azurewebsites.net/")

	if skipToken != "" {
		q := req.URL.Query()
		q.Add("$skiptoken", skipToken)
		req.URL.RawQuery = q.Encode()
	}

	fmt.Println("Requesting with ", req.URL.String())
	usersResp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(usersResp.Body)
	if err != nil {
		log.Fatal(err)
	}
	adResp := &AzureADResponse{}
	err = json.Unmarshal(body, adResp)
	if err != nil {
		log.Fatal(err)
	}
	return adResp
}

func (p *Populator) getToken() *MSTokenResponse {
	data := url.Values{}
	data.Set("client_id", p.clientID)
	data.Set("client_secret", p.clientSecret)
	data.Set("grant_type", "client_credentials")
	data.Set("resource", "https://graph.windows.net")
	resp, err := http.Post("https://login.microsoftonline.com/solarwinds.com/oauth2/token", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	response := &MSTokenResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
