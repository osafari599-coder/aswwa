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
	"io/ioutil" // برای خواندن پاسخ سرور لایسنس
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
// آدرس دقیق فایل متنی در گیت‌هاب خودت را اینجا بگذار
const LicenseServerURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"

func getMachineID() string {
	id, _ := os.Hostname()
	return strings.TrimSpace(id)
}

func verifyLicenseOnline() bool {
	client := http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(LicenseServerURL)
	if err != nil {
		fmt.Println("❌ Error: Could not connect to license server.")
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	allowedList := string(body)
	mID := getMachineID()

	// بررسی خط به خط برای پیدا کردن Machine ID
	lines := strings.Split(allowedList, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == mID {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------

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

// ... بقیه توابع کمکی (rateLimitedConn و غیره) که در کد شما بود ...

func main() {
	// --- مرحله اول: تایید لایسنس ---
	fmt.Println("🔍 Verifying License...")
	if !verifyLicenseOnline() {
		fmt.Println("\n\033[31m=======================================")
		fmt.Println(" ❌ ACCESS DENIED: UNAUTHORIZED SERVER")
		fmt.Printf(" Your Machine ID: %s\n", getMachineID())
		fmt.Println(" Please contact Admin to whitelist this ID.")
		fmt.Println("=======================================\033[0m")
		os.Exit(1)
	}
	fmt.Println("✅ License Verified.")

	// --- ادامه اجرای برنامه اصلی ---
	mode := flag.String("mode", "", "internal: 'server' or 'client'")
	rateLimit := flag.Int("ratelimit", 0, "Max bytes per second per conn")
	dashboardPort := flag.String("dashboard", "", "Dashboard port")
	tunnelType := flag.String("tunnel-type", "wss", "Tunnel protocol")
	authToken := flag.String("token", "", "Authentication token")
	fragSize := flag.Int("frag-size", 0, "Fragmentation size")
	fragDelay := flag.Int("frag-delay", 0, "Fragmentation delay")
	flag.Parse()

	if *mode != "" {
		// کدهای مربوط به اجرای مخفی در پس‌زمینه
		// ... (همان کدهایی که فرستاده بودید)
	}
	showInteractiveMenu()
}

// سایر توابع (showInteractiveMenu, setupServer, runServer و غیره)
// را دقیقاً طبق فایل قبلی خودتان در اینجا قرار دهید...
