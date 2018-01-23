# gitlic-check
A small program to print all repository names and their associated license type for each organization you belong to.

## Usage
You'll need to generate personal access token [in your settings](https://github.com/settings/applications#personal-access-tokens). Put your username and personal access token, separated by a comma, in "user.txt", which will be read into the tool. Example:

`maxgardner,PAT`

### Sample Output
```
Org: cncf
  Repo: demo [Apache License 2.0]
  Repo: landscape [Apache License 2.0]
```