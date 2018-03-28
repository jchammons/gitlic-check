# gitlic-check
A small program to generate CSV reports about all of the organizations you own and the repositories and users under those organizations.

## Configuration
The program requires a `config.json` file to run. The options for this config file are:

```
github:
  pat: string
  ignoredOrgs: []string (optional)
  rmInvitesAfter: int (optional) // in hours (ex: 2 weeks = 336)
drive: (optional)
  outputDir: string // id for output directory
  enableTeamDrive: bool (optional)
```

#### GitHub
You'll need to generate a personal access token [in your GitHub settings](https://github.com/settings/applications#personal-access-tokens). If you want to include private repositories in your reports, be sure to select the entire repo scope in the token settings. If you want to ignore any orgs in the process of scanning, put their names in the optional array. Finally, if you want to automatically remove invitations after a certain amount of time, include that option and the time frame in hours.

#### Google Drive (optional)
If you'd like to upload the resulting CSVs directly into a Google Drive folder, you'll need to create a service account in the [Google Developer Console](https://console.developers.google.com/apis/) and enable access to the Google Drive API. Place the JSON key file they provide in this folder and rename it `config-drive.json`.

Then, include the _drive_ property in your config file with the ID of the output Google Drive folder. **Before running, you must share this folder with the email address of the service account**.

If this folder is in a Team Drive, you must add the service account as a member on the Drive and enable this option in your `config.json`.

## Usage
After adding your config.json file and your config-drive.json file (if you're uploading the reports to a Drive folder), run:
```
go run *.go
```
This will perform a full check, going through all of your organizations and their repositories, users, and invitations, depending on what configuration options you've enabled.

#### Optional Flags
There are some optional flags if you want to quickly omit a step without altering your config file.

_To skip the GitHub scan and test only the upload_
```
go run *.go -upload-only
```
_To skip the upload step and perform only the GitHub scan_
```
go run *.go -no-upload
```

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
Org,User,Date Sent,Invited By,Deleted
org-name,user-name,2018-03-15,user-name,true
```