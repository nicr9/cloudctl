package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"os"
	"strconv"
	"strings"
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
	fmt.Fprintln(w, "DropletID\tName\tPublic IP\tPrivateIP")
	fmt.Fprintln(w, "---\t---\t---\t---")

	total := 0
	for _, drop := range list {
		// Get lists for public and private networks
		nets := drop.Networks.V4
		var public []string
		var private []string
		for _, net := range nets {
			if net.Type == "public" {
				public = append(public, net.IPAddress)
			} else if net.Type == "private" {
				private = append(private, net.IPAddress)
			}
		}

		// Print instance details using tabwriter
		fmt.Fprintf(
			w,
			"%d\t%s\t%s\t%s\n",
			drop.ID,
			drop.Name,
			strings.Join(public, ","),
			strings.Join(private, ","),
		)
		total++
	}
	w.Flush()
	fmt.Printf("---\nFound %d instances.\n", total)
}

func (d DigitalOcean) showInstance(instanceId string) {
	inst := d.getDroplet(instanceId)
	if inst != nil {
		fmt.Printf("%#v\n", inst)
	} else {
		fmt.Printf("Couldn't find %s\n", instanceId)
	}
}

func (d DigitalOcean) sshInstance(username, instanceId string) {
}

func (d DigitalOcean) getDroplet(dropletId string) *godo.Droplet {
	id, err := strconv.ParseInt(dropletId, 10, 0)
	if err != nil {
		return nil
	}

	result, _, err := d.svc.Droplets.Get(int(id))
	if err != nil {
		return nil
	}

	return result
}
