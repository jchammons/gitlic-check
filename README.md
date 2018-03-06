# gitlic-check
A small program to generate CSV reports on the GitHub repositories for organizations of which you are a member. What it currently pulls:

## Usage
You'll need to generate a personal access token [in your settings](https://github.com/settings/applications#personal-access-tokens). If you want to include private repositories in your reports, be sure to select the entire repo scope in the token settings. Then, put your PAT in a file named "auth.txt", which will be read into the tool.

### Ignoring repositories
To exclude specific repositories from your reports, include their names separated by a comma in a file called "ignore.txt", which will be read into the tool.

### Sample Output
_repos.csv_
```
Org,Repo,Private,Fork,License 
org-name,repo-name,true,false,MIT License
```

_users.csv_
>Note: Must be an owner of the org to pull data on private members and two-factor authentication
```
Org,User,Two-Factor Enabled
org-name,user-name,true
```
