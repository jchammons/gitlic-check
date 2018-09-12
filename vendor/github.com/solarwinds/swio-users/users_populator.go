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

var (
	USER_TYPE  = "#microsoft.graph.user"
	GROUP_TYPE = "#microsoft.graph.group"
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

type ADResponse struct {
	Value    []ADObject `json:"value"`
	NextLink string     `json:"@odata.nextLink"`
}

type ADObject struct {
	DisplayName string `json:"displayName"`
	Enabled     bool   `json:"accountEnabled"`
	GivenName   string `json:"givenName"`
	Mail        string `json:"mail"`
	ObjectID    string `json:"id"`
	ObjectType  string `json:"@odata.type"`
	Surname     string `json:"surname"`
}

type Populator struct {
	clientID        string
	clientSecret    string
	engineeringOnly bool
	groupIDs        []string
	more            bool
	moreDisabled    bool
	nextLink        string
	token           *MSTokenResponse
	userCount       int
}

func NewPopulator(clientID, clientSecret string, topLevelGroups []string) *Populator {
	return &Populator{
		clientID:     clientID,
		clientSecret: clientSecret,
		groupIDs:     topLevelGroups,
		more:         true,
		moreDisabled: true,
	}
}

func (p *Populator) MoreUsers() bool {
	return p.more
}

func (p *Populator) MoreDisabled() bool {
	return p.moreDisabled
}

func (p *Populator) MoreGroups() bool {
	return len(p.groupIDs) > 0
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

func (p *Populator) GetDisabledUsers() ([]*User, error) {
	token := p.getToken()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users?$filter=accountEnabled%20eq%20false", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Referer", "https://graphexplorer.azurewebsites.net/")

	fmt.Println("Requesting with ", req.URL.String())
	usersResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(usersResp.Body)
	if err != nil {
		return nil, err
	}
	adResp := &ADResponse{}
	err = json.Unmarshal(body, adResp)
	if err != nil {
		return nil, err
	}

	users := []*User{}
	for _, val := range adResp.Value {
		users = append(users, &User{
			FirstName: val.GivenName,
			LastName:  val.Surname,
			Email:     val.Mail,
			Enabled:   val.Enabled,
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
		p.moreDisabled = false
	}
	return users, nil
}

func (p *Populator) GetUsersByGroup() ([]*User, error) {
	toke := p.getToken()
	if len(p.groupIDs) > 0 {
		id := p.groupIDs[0]
		p.groupIDs = p.groupIDs[1:]
		users, newGroups, err := p.getGroupMembers(toke, id)
		if err != nil {
			return nil, err
		}
		p.groupIDs = append(p.groupIDs, newGroups...)
		return users, nil
	}
	return []*User{}, nil
}

func (p *Populator) getGroupMembers(token *MSTokenResponse, groupId string) ([]*User, []string, error) {
	newGroups := []string{}
	users := []*User{}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://graph.microsoft.com/v1.0/groups/%s/members?$top=999&$select=displayName,accountEnabled,givenName,mail,objectId,objectType,surname", groupId), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	groupResp, err := client.Do(req)
	if err != nil {
		return nil, newGroups, err
	}
	body, err := ioutil.ReadAll(groupResp.Body)
	if err != nil {
		return nil, newGroups, err
	}

	adResp := &ADResponse{}
	err = json.Unmarshal(body, adResp)
	if err != nil {
		return nil, newGroups, err
	}
	for _, val := range adResp.Value {
		if val.ObjectType == GROUP_TYPE {
			newGroups = append(newGroups, val.ObjectID)
		} else if val.ObjectType == USER_TYPE {
			if val.Enabled {
				users = append(users, &User{
					FirstName: val.GivenName,
					LastName:  val.Surname,
					Email:     val.Mail,
					Enabled:   true,
				})
			}
		}
	}
	return users, newGroups, nil
}

func (p *Populator) requestUsers(token *MSTokenResponse, skipToken string) (*ADResponse, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users?$top=999", nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	if skipToken != "" {
		q := req.URL.Query()
		decodedToken, _ := url.QueryUnescape(skipToken)
		q.Add("$skiptoken", decodedToken)
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

	adResp := &ADResponse{}
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
	data.Set("resource", "https://graph.microsoft.com")
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
