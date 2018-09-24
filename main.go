package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("please set GITHUB_TOKEN env var with a Personal Access Token")
	}

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("please provide organisation name")
	}

	org := os.Args[1]

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	rch := make(chan []*github.Repository, 10)

	go func() {
		for {
			select {
			case repos := <-rch:
				for _, r := range repos {
					rurl := strings.TrimPrefix(*r.CloneURL, "https://")
					rurl = strings.TrimSuffix(rurl, ".git")
					fmt.Println(rurl)
				}
			}
		}
	}()

	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			log.Fatal(err)
		}
		rch <- repos
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

}
