# What Is This?
This repo contains a tool to manage SolarWinds-related source code management concerns. The various areas
of concern are split into separate subcommands off of the root, each of which will be detailed below.

## Building
The tool is written in Go and requires a local Go setup to successfully build and run locally. Simply running
`go build` will give you a binary that can then be used as described below. The binary will be called `augit` for
purposes of documentation but will vary depending on the name of your local directory.

## GitLic
A utility to verify and manage OSS licenses, contribution guidelines, and GitHub organization and repository settings for both public and private repositories.

### Functionality
*Audit* 
- Verify GitHub Organization and Repository settings (specified either globally or per org/repo).
  * Including two-factor authentication, presence, and type of LICENSE and CONTRIBUTING files.
  * and other org/repo specific settings.
- Catalog unaccepted invitations to repositories.

*Enforce*
- Apply GitHub Organization and Repository settings (specified either globally or per org/repo).
  * Including requiring two-factor authentication and other org/repo specific settings.
- Cancel unaccepted invitations of a certain age (e.g. 14 days old) to organizations.

*Report*
- Generate CSV reports on the GitHub repositories for organizations of which you are a member.
- Upload reports to Google Sheets.

## Configuration
GitLic requires a `config.json` file to run. The options for this config file are:

```
github:
  pat: string
  ignoredOrgs: []string (optional)
  rmInvitesAfter: int (optional) // in hours (ex: 2 weeks = 336)
drive: (optional)
  outputDir: string // id for output directory
  enableTeamDrive: bool (optional)
```

##### GitHub
You'll need to generate a personal access token [in your GitHub settings](https://github.com/settings/applications#personal-access-tokens). If you want to include private repositories in your reports, be sure to select the entire repo scope in the token settings. If you want to ignore any orgs in the process of scanning, put their names in the optional array. Finally, if you want to automatically remove invitations after a certain amount of time, include that option and the time frame in hours.

##### Google Drive (optional)
If you'd like to upload the resulting CSVs directly into a Google Drive folder, you'll need to create a service account in the [Google Developer Console](https://console.developers.google.com/apis/) and enable access to the Google Drive API. Place the JSON key file they provide in this folder and rename it `config-drive.json`.

Then, include the _drive_ property in your config file with the ID of the output Google Drive folder. **Before running, you must share this folder with the email address of the service account**.

If this folder is in a Team Drive, you must add the service account as a member on the Drive and enable this option in your `config.json`.

### Usage
After adding your config.json file and your config-drive.json file (if you're uploading the reports to a Drive folder), run:
```
./augit gitlic
```
This will perform a full check, going through all of your organizations and their repositories, users, and invitations, depending on what configuration options you've enabled.

##### Optional Flags
There are some optional flags if you want to quickly omit a step without altering your config file.

_To skip the GitHub scan and test only the upload_
```
./augit gitlic -upload-only
```
_To skip the upload step and perform only the GitHub scan_
```
./augit gitlic -no-upload
```

### Sample Output
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