# gitlic-check
A small program to generate CSV reports on the GitHub repositories for organizations of which you are a member.

## Usage
You'll need to generate a personal access token [in your settings](https://github.com/settings/applications#personal-access-tokens). If you want to include private repositories in your reports, be sure to select the entire repo scope in the token settings. Then, put your username and personal access token, separated by a comma, in "auth.txt", which will be read into the tool. Example:

`maxgardner,PAT`

### Ignoring repositories
If you know specific repositories you aren't interested in seeing data from, include their names separated by a comma in a file called "ignore.txt", which will be read into the tool.

### Sample Output
_repos.csv_
```
Org,Repo,License 
cncf,demo,Apache License 2.0
cncf,landscape,Apache License 2.0
```

_users.csv_
>Note: Must be an owner of the org to pull data on private members and two-factor authentication
```
Org,User,Two-Factor Enabled?
cncf,leecalcote,true
solarwinds,maxgardner,true
```
