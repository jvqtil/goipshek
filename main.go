package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fatih/color"
)

type IpInfo struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Region   string `json:"region"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
	Org      string `json:"org"`
	AS       string `json:"as"`
}

func main() {
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	pink := color.New(color.FgMagenta).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Fatalf("Error getting IP: %v", err)
	}
	defer resp.Body.Close()

	var ipResponse map[string]string
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ipResponse)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	publicIP := ipResponse["ip"]
	fmt.Println()
	fmt.Printf("%s %s\n",
	blue(fmt.Sprintf("%-10s", "Public IP ")),
	wrap(bold(publicIP)))

	ipInfoResp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", publicIP))
	if err != nil {
		log.Fatalf("Error getting info by IP: %v", err)
	}
	defer ipInfoResp.Body.Close()

	var ipInfoNow IpInfo
	body, _ = ioutil.ReadAll(ipInfoResp.Body)
	err = json.Unmarshal(body, &ipInfoNow)
	if err != nil {
		log.Fatalf("Error parsing IP data: %v", err)
	}

	// Print all output
	fmt.Println()

	fmt.Printf("%s %s\n", 
	red(fmt.Sprintf("%-10s", "Country ")), 
	wrap(ipInfoNow.Country))

	fmt.Printf("%s %s\n", 
	green(fmt.Sprintf("%-10s", "Region ")), 
	wrap(ipInfoNow.Region))

	fmt.Printf("%s %s\n", 
	blue(fmt.Sprintf("%-10s", "City ")), 
	wrap(ipInfoNow.City))

	fmt.Printf("%s %s\n", 
	yellow(fmt.Sprintf("%-10s", "Provider ")), 
	wrap(ipInfoNow.ISP))

	fmt.Printf("%s %s\n", 
	pink(fmt.Sprintf("%-10s", "Org ")), 
	wrap(ipInfoNow.Org))

	fmt.Printf("%s %s\n", 
	cyan(fmt.Sprintf("%-10s", "AS ")), 
	wrap(ipInfoNow.AS))
}

func wrap(arg string) string {
	return "[ " + arg + " ]"
}
