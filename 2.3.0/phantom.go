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

// --- بخش لایسنس آنلاین ---
const LicenseURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"

func getMachineID() string {
	hostname, _ := os.Hostname()
	return strings.TrimSpace(hostname)
}

func verifyLicense() bool {
	client := http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(LicenseURL)
	if err != nil {
		fmt.Printf("❌ Connection Error: %v\n", err)
		return false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	mID := getMachineID()
	return strings.Contains(string(body), mID)
}

// --- تنظیمات ثابت فانتوم ---
const (
	logFilePath       = "/tmp/phantom-tunnel.log"
	pidFilePath       = "/tmp/phantom.pid"
	successSignalPath = "/tmp/phantom_success.signal"
)

var bufferPool = &sync.Pool{
	New: func() any { return make([]byte, 32*1024) },
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

// ... (توابع کمکی مثل rateLimitedConn و pipeCount و غیره را از کد اصلی خودت اینجا داشته باش)

func main() {
	// ۱. بررسی لایسنس قبل از منو
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

	// ۲. مدیریت فلگ‌ها و مد اجرای پس‌زمینه
	mode := flag.String("mode", "", "internal: 'server' or 'client'")
	rateLimit := flag.Int("ratelimit", 0, "Max bytes per second")
	dashboardPort := flag.String("dashboard", "", "Dashboard port")
	tunnelType := flag.String("tunnel-type", "wss", "Tunnel protocol")
	authToken := flag.String("token", "", "Auth token")
	fragSize := flag.Int("frag-size", 0, "Frag size")
	fragDelay := flag.Int("frag-delay", 0, "Frag delay")
	flag.Parse()

	if *mode != "" {
		// بخش اجرای سرویس در بک‌گراند (کدهایی که فرستاده بودی)
		configureLogging()
		dbPort := *dashboardPort
		if dbPort == "" {
			if *mode == "server" { dbPort = "8080" } else { dbPort = "8081" }
		}
		go startWebDashboard(":" + dbPort)
		
		args := flag.Args()
		if *mode == "server" && len(args) >= 5 {
			runServer(args[0], args[1], args[2], args[3], args[4], *rateLimit, *tunnelType, *authToken, *fragSize, *fragDelay)
		} else if *mode == "client" && len(args) >= 2 {
			runClient(args[0], args[1], *rateLimit, *tunnelType, *authToken, *fragSize, *fragDelay)
		}
		return
	}

	// ۳. نمایش منوی تعاملی برای کاربر مجاز
	showInteractiveMenu()
}

// ... بقیه توابع (showInteractiveMenu, runServer, runClient و غیره) بدون تغییر ...
