# gitlic-check
A utility to verify and manage OSS licenses, contribution guidelines, Github organization and repository settings for both public and private repositories.

## Functionality 
*Audit* 
- Verify Github Organization and Repository settings (specified either globally or per org/repo).
  * Including two-factor authentication, presence and type of LICENSE and CONTRIBUTING files.
  * and other org/repo specific settings.
- Catalog unaccepted invitations to repositories.

*Enforce*
- Apply Github Organization and Repository settings (specified either globally or per org/repo).
  * Including requiring two-factor authentication and other org/repo specific settings.
- Cancel unaccepted invitations of a certain age (e.g. 14 days old) to repositories.

*Report*
- Generate CSV reports on the GitHub repositories for organizations of which you are a member.
- Upload reports to Google Sheets.

## Usage
To connect with the GitHub API, generate a personal access token [in your settings](https://github.com/settings/applications#personal-access-tokens). If you want to include private repositories in your reports, be sure to select the entire repo scope in the token settings. Then, put your PAT in a file named `tokens.txt`, which will be read into the tool.

### Ignoring repositories
To exclude specific repositories from your reports, include their names separated by a comma in a file called `ignore.txt`, which will be read into the tool.

## Sample Output
>Note: Must be an owner of the org to pull data on private repos, members, and two-factor authentication

_repos.csv_
```
Org,Repo,Private,Fork,License 
org-name,repo-name,true,false,MIT License
```

_users.csv_
```
Org,User,Two-Factor Enabled
org-name,user-name,true
```

_invites.csv_
```
Org,User,Date Sent,Invited By
org-name,user-name,2018-03-15,user-name
```

## Uploading to Google Drive
If you'd like to upload the resulting CSVs directly into a Google Drive, you'll need to create a service account in the [Google Developer Console](https://console.developers.google.com/apis/) and enable access to the Google Drive API. Place the JSON key file in this folder and rename it `config.json`.

Next, put the ID of the Google Drive folder you'd like them uploaded to in `tokens.txt`, separated from your PAT by a comma. **Be sure to share the folder with the email address of the service account**.

If this folder is in a Team Drive, then you must add the service account as a member on the Drive and include the text "true" after your folder ID in `tokens.txt`.

### tokens.txt Sample
>Note: You do not need to specify false if it is not a Team Drive
```
github-PAT,drive-folder-ID,true
```
