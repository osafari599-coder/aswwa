package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"nhooyr.io/websocket"
)

// --- تنظیمات لایسنس آنلاین ---
const (
	LicenseURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"
)

// دریافت Machine ID (نام هاست سرور)
func getMachineID() string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		return "unknown-device"
	}
	return strings.TrimSpace(hostname)
}

// بررسی آنلاین لایسنس از گیت‌هاب
func verifyLicense() bool {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(LicenseURL)
	if err != nil {
		fmt.Printf("❌ Error: Could not connect to license server (%v)\n", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	allowedIDs := string(body)
	mID := getMachineID()

	// چک کردن اینکه آیا نام این سرور در فایل متنی گیت‌هاب شما هست یا نه
	lines := strings.Split(allowedIDs, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == mID {
			return true
		}
	}
	return false
}

// --- متغیرهای سیستمی فانتوم ---
const (
	logFilePath       = "/tmp/phantom-tunnel.log"
	pidFilePath       = "/tmp/phantom.pid"
	successSignalPath = "/tmp/phantom_success.signal"
)

var bufferPool = &sync.Pool{
	New: func() any {
		return make([]byte, 32*1024)
	},
}

type TunnelStats struct {
	sync.Mutex
	ActiveConnections int
	TotalBytesIn      int64
	TotalBytesOut     int64
	Uptime            time.Time
	Connected         bool
}

var stats = &TunnelStats{Uptime: time.Now()}

type activeSession struct {
	sync.RWMutex
	session *yamux.Session
}

func (as *activeSession) Get() *yamux.Session {
	as.RLock()
	defer as.RUnlock()
	return as.session
}

func (as *activeSession) Set(session *yamux.Session) {
	as.Lock()
	defer as.Unlock()
	if as.session != nil && !as.session.IsClosed() {
		as.session.Close()
	}
	as.session = session
}

// --- شروع اجرای اصلی برنامه ---
func main() {
	// قدم اول: چک کردن لایسنس
	fmt.Println("🔍 Verifying License...")
	if !verifyLicense() {
		fmt.Println("\n\033[31m##########################################")
		fmt.Println("       LICENSE ERROR: UNAUTHORIZED")
		fmt.Printf("       Your Machine ID: %s\n", getMachineID())
		fmt.Println("   Contact Admin to whitelist your server.")
		fmt.Println("##########################################\033[0m")
		os.Exit(1)
	}
	fmt.Println("✅ License Verified.")

	// تنظیمات ورودی (Flags)
	mode := flag.String("mode", "", "internal: 'server' or 'client'")
	rateLimit := flag.Int("ratelimit", 0, "Max bytes per second per conn")
	dashboardPort := flag.String("dashboard", "", "Dashboard port")
	tunnelType := flag.String("tunnel-type", "wss", "Tunnel protocol")
	authToken := flag.String("token", "", "Authentication token")
	fragSize := flag.Int("frag-size", 0, "Fragmentation size")
	fragDelay := flag.Int("frag-delay", 0, "Fragmentation delay")
	flag.Parse()

	if *mode != "" {
		configureLogging()
		args := flag.Args()
		dbPort := *dashboardPort
		if dbPort == "" {
			dbPort = "8080"
			if *mode == "client" { dbPort = "8081" }
		}
		go startWebDashboard(":" + dbPort)
		
		if *mode == "server" {
			if len(args) < 5 { log.Fatal("Missing server arguments") }
			runServer(args[0], args[1], args[2], args[3], args[4], *rateLimit, *tunnelType, *authToken, *fragSize, *fragDelay)
		} else if *mode == "client" {
			if len(args) < 2 { log.Fatal("Missing client arguments") }
			runClient(args[0], args[1], *rateLimit, *tunnelType, *authToken, *fragSize, *fragDelay)
		}
		return
	}
	showInteractiveMenu()
}

// ... بقیه توابع شما (runServer, runClient, startWebDashboard و غیره) را بدون تغییر در ادامه کپی کنید ...
