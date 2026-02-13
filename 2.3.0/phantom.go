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
const (
	LicenseURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"
)

func getMachineID() string {
	hostname, _ := os.Hostname()
	return strings.TrimSpace(hostname)
}

func verifyLicense() bool {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(LicenseURL)
	if err != nil {
		fmt.Printf("❌ Error connecting to license server: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	mID := getMachineID()
	
	// بررسی وجود Hostname در فایل گیت‌هاب
	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == mID {
			return true
		}
	}
	return false
}

// --- متغیرها و ساختارهای اصلی فانتوم ---
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

// --- تابع اصلی (Main) ---
func main() {
	// ۱. بررسی لایسنس (اجباری)
	fmt.Println("🔍 Checking License...")
	if !verifyLicense() {
		fmt.Println("\n\033[31m##########################################")
		fmt.Println("       LICENSE ERROR: UNAUTHORIZED")
		fmt.Printf("       Your Machine ID: %s\n", getMachineID())
		fmt.Println("   Contact Admin to whitelist your server.")
		fmt.Println("##########################################\033[0m")
		os.Exit(1)
	}
	fmt.Println("✅ License Verified Successfully.")

	// ۲. تعریف فلگ‌ها برای تنظیمات و اجرا
	mode := flag.String("mode", "", "internal: 'server' or 'client'")
	setupPort := flag.String("setup-port", "", "Port for initial setup")
	setupUser := flag.String("setup-user", "", "User for initial setup")
	setupPass := flag.String("setup-pass", "", "Pass for initial setup")
	startPanel := flag.Bool("start-panel", false, "Start the web dashboard")
	
	// سایر فلگ‌های مربوط به تونل
	rateLimit := flag.Int("ratelimit", 0, "Max bytes per second")
	tunnelType := flag.String("tunnel-type", "wss", "Tunnel protocol")
	flag.Parse()

	// ۳. مدیریت بخش Setup (جلوگیری از ارور Too few arguments)
	if *setupPort != "" && *setupUser != "" && *setupPass != "" {
		fmt.Printf("⚙️ Configuring Phantom on port %s...\n", *setupPort)
		// در اینجا کد ذخیره تنظیمات در دیتابیس یا فایل را قرار بده
		// فعلاً یک فایل سیگنال برای اتمام نصب می‌سازیم
		os.WriteFile(successSignalPath, []byte("ok"), 0644)
		fmt.Println("✅ Setup completed.")
		return
	}

	// ۴. اجرای پنل یا منوی اصلی
	if *startPanel {
		fmt.Println("🚀 Starting Web Dashboard...")
		// کدهای مربوط به startWebDashboard را اینجا فراخوانی کن
		select {} // نگه داشتن برنامه
	}

	// ۵. اگر هیچ آرگومانی نبود، منوی تعاملی باز شود
	showInteractiveMenu()
}

func showInteractiveMenu() {
	fmt.Println("\n--- Phantom Tunnel Interactive Menu ---")
	fmt.Println("1. Start Server")
	fmt.Println("2. Start Client")
	fmt.Println("3. Exit")
	// کدهای منوی خودت را اینجا ادامه بده...
}

// بقیه توابع شما (runServer, runClient, غیره) را در ادامه کپی کنید...
