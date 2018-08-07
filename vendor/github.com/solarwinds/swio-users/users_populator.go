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
	Enabled   bool
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
	Enabled   bool   `json:"accountEnabled"`
}

type Populator struct {
	clientID     string
	clientSecret string
	more         bool
	token        *MSTokenResponse
	nextLink     string
}

func NewPopulator(clientID, clientSecret string) *Populator {
	return &Populator{
		clientID:     clientID,
		clientSecret: clientSecret,
		more:         true,
	}
}

func (p *Populator) MoreUsers() bool {
	return p.more
}

func (p *Populator) GetUsers() ([]*User, error) {
	toke := p.getToken()
	adResp, err := p.requestUsers(toke, p.nextLink)
	if err != nil {
		return nil, err
	}
	users := []*User{}
	for _, user := range adResp.Value {
		if user.Mail == "" {
			fmt.Printf("[NO EMAIL] skipping user: %+v\n", user)
			continue
		}
		users = append(users, &User{
			FirstName: user.GivenName,
			LastName:  user.Surname,
			Email:     user.Mail,
			Enabled:   user.Enabled,
		})
	}

	nextLinkReg := regexp.MustCompile("skiptoken=(.*)")
	matches := nextLinkReg.FindStringSubmatch(adResp.NextLink)
	if len(matches) == 2 {
		p.nextLink = matches[1]
	} else {
		p.nextLink = ""
	}
	if adResp.NextLink == "" {
		p.more = false
	}
	return users, nil
}

func (p *Populator) requestUsers(token *MSTokenResponse, skipToken string) (*AzureADResponse, error) {
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
		return nil, err
	}
	body, err := ioutil.ReadAll(usersResp.Body)
	if err != nil {
		return nil, err
	}
	adResp := &AzureADResponse{}
	err = json.Unmarshal(body, adResp)
	if err != nil {
		return nil, err
	}
	return adResp, nil
}

func (p *Populator) getToken() *MSTokenResponse {
	if p.token != nil {
		return p.token
	}
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
	p.token = response
	return response
}
