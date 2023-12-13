package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/transip/gotransip/v6"
	"github.com/transip/gotransip/v6/domain"
)

func getPublicIP() string {
	response, err := http.Get("https://ipecho.net/plain")
	if err != err {
		return fmt.Sprintln("Couldn't query public IP")
	} else {
		defer response.Body.Close()
		resBody, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		bodyString := string(resBody)
		return bodyString
	}
}

func updateRootDNSEntry(publicIP string) {
	publicIPAddress := publicIP
	fmt.Println("Public IP address is: ", publicIPAddress)
	// create a new TransIP API client
	// You need to create an API key in the Portal first,
	// make sure IP whitelisting is not enabled, won't work for this purpose
	client, err := gotransip.NewClient(gotransip.ClientConfiguration{
		AccountName:    "tranip_username_here", //transIP username here
		PrivateKeyPath: "./private.key",        //path to private key
	})
	if err != nil {
		panic(err.Error())
	}
	domainRepo := domain.Repository{Client: client}

	// get a list of your Domains
	dns, err := domainRepo.GetAll()
	if err != nil {
		panic(err.Error())
	}

	selectedDns := dns[0].Name
	fmt.Println("Selected domainname is: ", selectedDns)
	dnsEntries, err := domainRepo.GetDNSEntries(selectedDns)
	var currentRootEntry domain.DNSEntry
	for _, dnsentry := range dnsEntries {
		if dnsentry.Name == "@" && dnsentry.Type == "A" {
			currentRootEntry = dnsentry
			fmt.Println("Root DNS entry selected: ", dnsentry)
		}
	}
	if currentRootEntry.Content != publicIPAddress {
		fmt.Println("Current public IP is different, udpating...")
		newRootEntry := currentRootEntry
		newRootEntry.Content = publicIPAddress
		domainRepo.UpdateDNSEntry(selectedDns, newRootEntry)
	} else {
		fmt.Println("Current public IP is equal to DNS root entry, exiting...")
	}
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	updateRootDNSEntry(getPublicIP())
}
