package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/fatih/color"
)

type ipInfo struct {
	IP      string `json:"query"`
	Country string `json:"country"`
	RG      string `json:"region"`
	Region  string `json:"regionName"`
	City    string `json:"city"`
	ISP     string `json:"isp"`
	Org     string `json:"org"`
	AS      string `json:"as"`
}

var (
	white  = color.New(color.FgWhite).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	pink   = color.New(color.FgMagenta).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
)

func main() {
	var toCheck string
	if len(os.Args) > 1 {
		toCheck = os.Args[1]
	} else {
		toCheck = getUserIP()
	}

	localip, err := getLocalIP()
	if err != nil {
		fmt.Println(err)
	}

	// Print all output
	fmt.Println()
	printer("Local IP", bold(localip), white)
	printer("Public IP", bold(toCheck), blue)

	ipInfoNow := getInfo(toCheck)

	if ipInfoNow.IP != toCheck {
		printer("Raw IP", bold(ipInfoNow.IP), yellow)
	}

	fmt.Println()

	printer("Country", ipInfoNow.Country, red)
	RegionData := ipInfoNow.RG + " | " + ipInfoNow.Region
	if ipInfoNow.RG != "" {
		printer("Region", RegionData, green)
	}
	printer("City", ipInfoNow.City, blue)
	printer("Provider", ipInfoNow.ISP, yellow)
	printer("Org", ipInfoNow.Org, pink)
	printer("AS", ipInfoNow.AS, cyan)
}

func wrap(arg string) string {
	return "[ " + arg + " ]"
}

func printer(label string, value string, colorFunc func(a ...interface{}) string) {
	if value != "" {
		fmt.Printf("%s %s\n", colorFunc(fmt.Sprintf("%-10s", label)), wrap(value))
	}
}

func getUserIP() string {
	resp, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		log.Fatalf("Error getting IP: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		IP string `json:"ip"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Error decoding IP: %v", err)
	}
	toCheck := result.IP

	return toCheck
}

func getInfo(toCheck string) ipInfo {
	resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", toCheck))
	if err != nil {
		log.Fatalf("Error getting info by IP: %v", err)
	}
	defer resp.Body.Close()

	var ipInfoNow ipInfo
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ipInfoNow)
	if err != nil {
		log.Fatalf("Error parsing IP data: %v", err)
	}

	return ipInfoNow
}

func getLocalIP() (string, error) {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}
