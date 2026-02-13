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

// تنظیمات لایسنس آنلاین
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
	mID := getMachineID()
	return strings.Contains(string(body), mID)
}

func main() {
	// چک کردن لایسنس
	if !verifyLicense() {
		fmt.Printf("\n\033[31m❌ ACCESS DENIED! Your Machine ID (%s) is not authorized.\033[0m\n", getMachineID())
		os.Exit(1)
	}

	// تعریف فلگ‌ها برای جلوگیری از ارور آرگومان
	setupPort := flag.String("setup-port", "", "Setup port")
	setupUser := flag.String("setup-user", "", "Setup user")
	setupPass := flag.String("setup-pass", "", "Setup pass")
	startPanel := flag.Bool("start-panel", false, "Start panel")
	flag.Parse()

	if *setupPort != "" {
		fmt.Printf("⚙️ Configuring Phantom on port %s...\n", *setupPort)
		return
	}

	if *startPanel {
		fmt.Println("🚀 Phantom Tunnel is running...")
		select {} 
	}

	fmt.Println("Welcome to Phantom Tunnel v2.3.0")
}
