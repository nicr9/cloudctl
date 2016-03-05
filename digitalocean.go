package main

import (
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"os/user"
	"path"
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

func (d DigitalOcean) listMachines() {
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
		public, private := getNetworks(&drop)

		// Print machine details using tabwriter
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
	fmt.Printf("---\nFound %d machines.\n", total)
}

func getNetworks(drop *godo.Droplet) (public, private []string) {
	// Get lists for public and private networks
	nets := drop.Networks.V4
	for _, net := range nets {
		if net.Type == "public" {
			public = append(public, net.IPAddress)
		} else if net.Type == "private" {
			private = append(private, net.IPAddress)
		}
	}

	return
}

func (d DigitalOcean) showMachine(machineId string) {
	inst := d.getDroplet(machineId)
	if inst != nil {
		fmt.Printf("%#v\n", inst)
	} else {
		fmt.Printf("Couldn't find %s\n", machineId)
	}
}

func (d DigitalOcean) sshMachine(username, machineId string) {
	drop := d.getDroplet(machineId)
	if drop != nil {
		ip, err := drop.PublicIPv4()
		if d.config.PrivateNetwork {
			ip, err = drop.PrivateIPv4()
		}
		if err != nil {
			fmt.Println("Can't find an IP for", machineId)
			return
		}
		userHost := fmt.Sprintf("%s@%s", username, ip)

		me, err := user.Current()
		if err != nil {
			fmt.Println("Can't determine username details:", err)
		}
		keyFile := fmt.Sprintf(".ssh/%s", d.config.KeyName)
		keyFile = path.Join(me.HomeDir, keyFile)

		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			fmt.Println("Can't find private key:", keyFile)
			return
		}

		cmd := exec.Command("ssh", "-i", keyFile, userHost)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Couldn't find %s\n", machineId)
	}
}

func (d DigitalOcean) removeMachines(machines []string) {
	for _, drop := range machines {
		id, err := strconv.ParseInt(drop, 10, 0)
		if err != nil {
			fmt.Println("Couldn't find", drop)
			continue
		}
		d.svc.Droplets.Delete(int(id))
	}
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
