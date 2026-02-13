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
	return strings.Contains(string(body), getMachineID())
}

// --- ساختارهای آماری ---
type TunnelStats struct {
	sync.Mutex
	ActiveConnections int
	TotalBytesIn      int64
	TotalBytesOut     int64
	Uptime            time.Time
	Connected         bool
}
var stats = &TunnelStats{Uptime: time.Now()}

// --- تابع اصلی ---
func main() {
	// ۱. چک کردن لایسنس در شروع برنامه
	if !verifyLicense() {
		fmt.Printf("\n\033[31m❌ Access Denied! Your Machine ID (%s) is not authorized.\033[0m\n", getMachineID())
		os.Exit(1)
	}

	// ۲. تعریف آرگومان‌ها (Flags)
	mode := flag.String("mode", "", "server or client")
	setupPort := flag.String("setup-port", "", "Port for setup")
	setupUser := flag.String("setup-user", "", "User for setup")
	setupPass := flag.String("setup-pass", "", "Pass for setup")
	flag.Parse()

	// ۳. اگر دستور ستاپ از سمت install.sh اومده باشه
	if *setupPort != "" {
		fmt.Printf("⚙️ Setting up Phantom on port %s...\n", *setupPort)
		// اینجا می‌تونی دیتابیس یا فایل تنظیمات رو بسازی
		os.WriteFile("/tmp/phantom_success.signal", []byte("ok"), 0644)
		return
	}

	// ۴. اگر مد سرور یا کلاینت انتخاب شده باشه
	if *mode != "" {
		fmt.Printf("🚀 Running in %s mode...\n", *mode)
		// فراخوانی توابع runServer یا runClient
		select {} 
	}

	// ۵. در غیر این صورت منوی گرافیکی/تعاملی
	showMenu()
}

func showMenu() {
	fmt.Println("=======================================")
	fmt.Println(" 👻 Phantom Tunnel v2.3 Online Edition")
	fmt.Println("=======================================")
	fmt.Println("1. Start Server")
	fmt.Println("2. Exit")
	// بقیه منوی خودت...
}
