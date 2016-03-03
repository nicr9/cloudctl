package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"os"
	"text/tabwriter"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type DigitalOcean struct {
	config Config
	svc    *godo.Client
}

func NewDigitalOcean(config Config) DigitalOcean {
	tokenSource := &TokenSource{
		AccessToken: config.AccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	svc := godo.NewClient(oauthClient)
	return DigitalOcean{config, svc}
}

func (d DigitalOcean) listInstances() {
	// Get droplet list
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := d.svc.Droplets.List(opt)
		if err != nil {
			return
		}
		for _, d := range droplets {
			list = append(list, d)
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return
		}
		opt.Page = page + 1
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "InstanceID\tPublic IP\tPrivateIP")
	fmt.Fprintln(w, "---\t---\t---")

	total := 0
	for _, drop := range list {
		// Replace public ip with "-" if instance doesn't have one
		fmt.Fprintf(
			w,
			"%s\t%s\t%s\n",
			drop.Name,
			"-",
			"-",
		)
		total++
	}
	w.Flush()
	fmt.Printf("---\nFound %d instances.\n", total)
}

func (d DigitalOcean) showInstance(instanceId string) {
}

func (d DigitalOcean) sshInstance(username, instanceId string) {
}

func (d DigitalOcean) getInstance(instanceId string) *godo.Droplet {
	return nil
}