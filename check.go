import "golang.org/x/oauth2"

func main() {
  ts := oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: "... your access token ..."},
  )
  tc := oauth2.NewClient(oauth2.NoContext, ts)

  client := github.NewClient(tc)

  // list all repositories for the authenticated user
  repos, _, err := client.Repositories.List("", nil)
}


import "github.com/google/go-github/github"
client := github.NewClient(nil)

// list all organizations for user "leecalcote"
orgs, _, err := client.Organizations.List("leecalcote", nil)

