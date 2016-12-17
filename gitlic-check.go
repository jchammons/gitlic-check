package main

import (
	
	"fmt"
	"io/ioutil"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)
// you need to generate personal access token at
// https://github.com/settings/applications#personal-access-tokens
// Put your personal access token in pat.txt

func setPersonalAccessToken() string {
	fh, err := ioutil.ReadFile("pat.txt")
	if err != nil {
		fmt.Print(err)
	}
	pat := string(fh)
	//fmt.Println("Retrived pat: ",pat)
	return pat
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
func getAll(){
	client := github.NewClient(nil)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg("github", opt)
		if err != nil {
			//return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			fmt.Printf("a page")
			break
			
		}
		opt.ListOptions.Page = resp.NextPage
	}
}
func main() {

	personalAccessToken := setPersonalAccessToken()
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)
	
	orgs, _, err := client.Organizations.List("leecalcote", nil)
	
	for _, org := range orgs {
		if err !=nil {
			fmt.Printf("Failed with %s\n", err)
		return
		}
		fmt.Printf("Org: %s \n", *org.Login)

	
		repos, _, err := client.Repositories.List(*org.Login, nil)			
		if err !=nil {
			fmt.Printf("Failed with %s\n", err)
			return
		}
		for _, repo := range repos {
			//fmt.Printf("Repo: %s \n", *repo.Owner.Login)
			lics, _, err := client.Repositories.License(*repo.Owner.Login, *repo.Name)
			if err == nil {
				fmt.Printf("  Repo: %s [%s] \n", *repo.Name, *lics.License.Name)
				
			}
		}
	}
	//getAll()
	// for _, lic := range lics {
    // 	fmt.Printf("Lic: %s\n", *lic.Name)
    // }
	//fmt.Printf("orgs:\n%s\n", string(orgs))

	// user, _, err := client.Users.Get("")
	// if err != nil {
	// 	fmt.Printf("client.Users.Get() faled with '%s'\n", err)
	// 	return
	// }
	// d, err := json.MarshalIndent(user, "", "  ")
	// if err != nil {
	// 	fmt.Printf("json.MarshlIndent() failed with %s\n", err)
	// 	return
	// }
	// fmt.Printf("User:\n%s\n", string(d))
	
	

	//func (s *RepositoriesService) License(owner, repo string) (*RepositoryLicense, *Response, error)
	
	//lics, _, err := client.Licenses.List()
}
