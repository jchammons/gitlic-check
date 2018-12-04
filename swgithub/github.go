package swgithub

import (
	"context"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/solarwinds/gitlic-check/config"
)

func GetSWOrgs(ctx context.Context, ghClient *github.Client, cf config.Config) ([]*github.Organization, error) {
	includedOrgs := getIncludedMap(cf.Github.IncludedOrgs)
	orgs := []*github.Organization{}
	lo := &github.ListOptions{PerPage: 100}
	for {
		partialOrgs, resp, err := ghClient.Organizations.List(ctx, "", lo)
		if err != nil {
			log.Fatalf("Organizations.List failed with %s\n", err)
			return nil, err
		}

		orgs = append(orgs, partialOrgs...)

		if resp.NextPage == 0 {
			lo.Page = 1
			break
		}
		lo.Page = resp.NextPage
	}
	validOrgs := []*github.Organization{}
	for _, org := range orgs {
		if includedOrgs == nil {
			validOrgs = orgs
			break
		}
		if _, ok := includedOrgs[strings.ToLower(*org.Login)]; ok {
			validOrgs = append(validOrgs, org)
		} else {
			log.Printf("Ignored %s\n", *org.Login)
		}
	}
	return validOrgs, nil
}

func GetOrgInvites(ctx context.Context, ghClient *github.Client, org *github.Organization) ([]*github.Invitation, error) {
	var invites []*github.Invitation
	opt := &github.ListOptions{PerPage: 100}
	for {
		partialInvites, resp, err := ghClient.Organizations.ListPendingOrgInvitations(ctx, *org.Login, opt)
		if err != nil {
			return nil, err
		}

		invites = append(invites, partialInvites...)

		if resp.NextPage == 0 {
			opt.Page = 1
			break
		}
		opt.Page = resp.NextPage
	}
	return invites, nil
}

func GetOrgMembers(ctx context.Context, ghClient *github.Client, org *github.Organization, opt *github.ListMembersOptions) ([]*github.User, error) {
	members := []*github.User{}
	for {
		partialMembers, resp, err := ghClient.Organizations.ListMembers(ctx, *org.Login, opt)
		if err != nil {
			log.Printf("Organizations.ListMembers, no filter, failed with %s\n", err)
			return nil, err
		}

		members = append(members, partialMembers...)

		if resp.NextPage == 0 {
			opt.Page = 1
			break
		}
		opt.Page = resp.NextPage
	}
	return members, nil
}

func GetOrgOwners(ctx context.Context, ghClient *github.Client, org *github.Organization) ([]*github.User, error) {
	lo := &github.ListOptions{PerPage: 100}
	memOpt := &github.ListMembersOptions{
		ListOptions: *lo,
		Role:        "admin",
	}
	members := []*github.User{}
	for {
		partialMembers, resp, err := ghClient.Organizations.ListMembers(ctx, *org.Login, memOpt)
		if err != nil {
			log.Printf("Organizations.ListMembers, only admin, failed with %s\n", err)
			return nil, err
		}

		members = append(members, partialMembers...)

		if resp.NextPage == 0 {
			memOpt.Page = 1
			break
		}
		memOpt.Page = resp.NextPage
	}
	return members, nil
}

func GetOrgRepositories(ctx context.Context, ghClient *github.Client, org *github.Organization) ([]*github.Repository, error) {
	var repos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		partialRepos, resp, err := ghClient.Repositories.ListByOrg(ctx, *org.Login, opt)
		if err != nil {
			return nil, err
		}

		repos = append(repos, partialRepos...)

		if resp.NextPage == 0 {
			opt.Page = 1
			break
		}
		opt.Page = resp.NextPage
	}

	return repos, nil
}

func GetCommitActivity(ctx context.Context, ghClient *github.Client, org string, repo string) ([]*github.WeeklyCommitActivity, error, bool) {
	act, resp, err := ghClient.Repositories.ListCommitActivity(ctx, org, repo)
	if resp.StatusCode == 202 {
		return nil, nil, true
	} else if err != nil && resp.StatusCode != 202 {
		return nil, err, true
	}
	return act, nil, false
}

func GetAdditionsDeletions(ctx context.Context, ghClient *github.Client, org string, repo string) ([]*github.WeeklyStats, error, bool) {
	mods, resp, err := ghClient.Repositories.ListCodeFrequency(ctx, org, repo)
	if resp.StatusCode == 202 {
		return nil, nil, true
	} else if err != nil && resp.StatusCode != 202 {
		return nil, err, true
	}
	return mods, nil, false
}
