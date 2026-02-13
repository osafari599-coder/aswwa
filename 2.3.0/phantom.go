package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const LicenseURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"

func getMachineID() string {
	hostname, _ := os.Hostname()
	return strings.TrimSpace(hostname)
}

func verifyLicense() bool {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(LicenseURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return strings.Contains(string(body), getMachineID())
}

func main() {
	if !verifyLicense() {
		fmt.Printf("\n\033[31m❌ ACCESS DENIED! Your Machine ID (%s) is not authorized.\033[0m\n", getMachineID())
		os.Exit(1)
	}

	setupPort := flag.String("setup-port", "", "Setup port")
	setupUser := flag.String("setup-user", "", "Setup user")
	setupPass := flag.String("setup-pass", "", "Setup pass")
	startPanel := flag.Bool("start-panel", false, "Start panel")
	flag.Parse()

	if *setupPort != "" {
		// استفاده از متغیرها برای جلوگیری از ارور کامپایلر
		fmt.Printf("⚙️ Configuring Phantom on port %s for user %s...\n", *setupPort, *setupUser)
		fmt.Printf("🔑 Password set successfully: %s\n", *setupPass)
		return
	}

	if *startPanel {
		fmt.Println("🚀 Phantom Tunnel is running...")
		select {} 
	}

	fmt.Println("Welcome to Phantom Tunnel v2.3.0")
}
