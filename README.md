# What Is This?
This repo contains a tool to manage SolarWinds-related source code management concerns. The various areas
of concern are split into separate subcommands off of the root, each of which will be detailed below.

## Building
The tool is written in Go and requires a local Go setup to successfully build and run locally. Simply running
`go build` will give you a binary that can then be used as described below. The binary will be called `augit` for
purposes of documentation but will vary depending on the name of your local directory.

------

## GitLic
A utility to verify and manage OSS licenses, contribution guidelines, and GitHub organization and repository settings for both public and private repositories.

### Functionality
<details>
<summary><em>Audit (expand)</em></summary>

- Verify GitHub Organization and Repository settings (specified either globally or per org/repo).
  * Including two-factor authentication, presence, and type of LICENSE and CONTRIBUTING files.
  * and other org/repo specific settings.
- Catalog unaccepted invitations to repositories.

</details>

<details> 
<summary><em>Enforce (expand)</em></summary>

- Apply GitHub Organization and Repository settings (specified either globally or per org/repo).
  * Including requiring two-factor authentication and other org/repo specific settings.
- Cancel unaccepted invitations of a certain age (e.g. 14 days old) to organizations.

</details>

<details>
<summary><em>Report (expand)</em></summary>

- Generate CSV reports on the GitHub repositories for organizations of which you are a member.
- Upload reports to Google Sheets.

</details>

### Usage
After adding your `options.json` file and, optionally, your `drive-key.json` file, run:
```
./augit gitlic
```
This will perform a full check, going through all of your organizations and their repositories, users, and invitations, depending on what configuration options you've enabled.

#### Optional Flags
_Skip the GitHub scan and test only the upload_
```
./gitlic-check gitlic --upload-only
```
_Skip the upload step and perform only the GitHub scan_
```
./gitlic-check gitlic --no-upload
```

## Configuration
GitLic requires a config file, located at `config/options.json`:

```js
{
  "github":
    "pat": "", // string
    "ignoredOrgs": [], // []string (optional)
    "includedOrgs": [], // []string (optional)
    "rmInvitesAfter": 336, // int (optional) // in hours (ex: 2 weeks = 336)
  "drive": // (optional)
    "outputDir": "", // string // id for output directory
    "enableTeamDrive": true // bool (optional)
}
```

#### GitHub
You'll need to generate a personal access token [in your GitHub settings](https://github.com/settings/tokens). If you want to include private repositories in your reports, be sure to select the entire repo scope in the token settings. If you want to ignore any orgs in the process of scanning, put their names in the optional array. Finally, if you want to automatically remove invitations after a certain amount of time, include that option and the time frame in hours.

The gh_report command populates the Augit database with GitHub users and instead of using the `ignoredOrgs` setting it uses `includedOrgs` as a whitelist.

>Note: Must be an owner of the org to pull data on private repos, members, and two-factor authentication

#### Google Drive (optional)
If you'd like to upload the resulting CSVs directly into a Google Drive folder, create
a service account in the [Google Developer Console](https://console.developers.google.com/apis/) and enable
access to the Google Drive API. Place the JSON key file they provide at `config/drive-key.json`; then, update the options file.

>Note: Before running, you must share the output folder with the service account's email address.
>If this folder is in a Team Drive, you must add the service account as a member on the Drive and enable this option in your `options.json`.

### Sample Output
<details>

<summary>Expand</summary>

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
</details>

-----

## Augit
Augit provides functionality to ensure that no GitHub account has access to SolarWinds repos unless they are an approved
and active employee or an approved bot account. This is accomplished by checking the list of active GitHub users in SolarWinds
organizations against a database of SolarWinds employees retrieved from Azure Active Directory. The Azure AD integration
is facilitated by the [swio-users](https://github.com/solarwinds/swio-users) repo, also used by Kudos and other
SolarWinds.io-related projects.

### Populating the Database
The `populate` subcommand will use your database.yml and env vars for ENVIRONMENT, DATABASE_URL (if production environment),
AD_CLIENT_ID, and AD_SECRET to populate your database with all enabled users from Active Directory. In the future it will also remove folks who are no longer enabled or in AD.

### Prerequisites for Completion
Augit relies on the personal access token of a user having *owner* level access to each SolarWinds repo. This is
required because the GitHub API will not return concealed users to any token that is not an owner.
