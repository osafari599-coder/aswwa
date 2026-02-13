package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"nhooyr.io/websocket"
)

// --- تنظیمات لایسنس آنلاین ---
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
	
	// بررسی دقیق نام سرور در فایل لیست سفید
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == mID {
			return true
		}
	}
	return false
}

// --- ساختارهای آماری و متغیرهای فانتوم ---
type TunnelStats struct {
	sync.Mutex
	ActiveConnections int
	TotalBytesIn      int64
	TotalBytesOut     int64
	Uptime            time.Time
	Connected         bool
}
var stats = &TunnelStats{Uptime: time.Now()}

// --- تابع اصلی اجرا ---
func main() {
	// ۱. تایید لایسنس قبل از هر عملیاتی
	if !verifyLicense() {
		fmt.Printf("\n\033[31m❌ ACCESS DENIED! Your Machine ID (%s) is not authorized.\033[0m\n", getMachineID())
		os.Exit(1)
	}

	// ۲. تعریف ورودی‌ها (برای جلوگیری از ارور Too few arguments)
	mode := flag.String("mode", "", "server or client")
	setupPort := flag.String("setup-port", "", "Setup port")
	setupUser := flag.String("setup-user", "", "Setup username")
	setupPass := flag.String("setup-pass", "", "Setup password")
	startPanel := flag.Bool("start-panel", false, "Start the panel service")
	flag.Parse()

	// ۳. مدیریت بخش ستاپ خودکار
	if *setupPort != "" {
		fmt.Printf("⚙️ Configuring Phantom on port %s...\n", *setupPort)
		// سیگنال موفقیت برای اسکریپت نصب
		os.WriteFile("/tmp/phantom_success.signal", []byte("ok"), 0644)
		return
	}

	// ۴. اجرای پنل وب
	if *startPanel {
		fmt.Println("🚀 Phantom Dashboard is starting...")
		// در اینجا تابع اجرای سرور وب خود را فراخوانی کنید
		select {} 
	}

	// ۵. اجرای مد سرور/کلاینت یا منوی اصلی
	if *mode == "server" {
		fmt.Println("Running in Server Mode...")
	} else if *mode == "client" {
		fmt.Println("Running in Client Mode...")
	} else {
		showMenu()
	}
}

func showMenu() {
	fmt.Println("\n=======================================")
	fmt.Println(" 👻 Phantom Tunnel v2.3.0 | Authorized")
	fmt.Println("=======================================")
	fmt.Println("1. Start Tunnel Server")
	fmt.Println("2. Start Tunnel Client")
	fmt.Println("3. Exit")
}
