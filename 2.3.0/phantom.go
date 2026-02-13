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
const LicenseURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"

func getMachineID() string {
	hostname, _ := os.Hostname()
	return strings.TrimSpace(hostname)
}

func verifyLicense() bool {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(LicenseURL)
	if err != nil {
		fmt.Println("❌ Error connecting to license server!")
		return false
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	allowedList := string(body)
	mID := getMachineID()

	// چک کردن وجود Machine ID در لیست مجاز
	lines := strings.Split(allowedList, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == mID {
			return true
		}
	}
	return false
}

// --- متغیرها و ساختارهای اصلی برنامه شما ---
const (
	logFilePath       = "/tmp/phantom-tunnel.log"
	pidFilePath       = "/tmp/phantom.pid"
	successSignalPath = "/tmp/phantom_success.signal"
)

// ... (بقیه متغیرها و توابع کمکی مثل TunnelStats و غیره که در کد خودت بود)

func main() {
	// ۱. تایید لایسنس قبل از هر کاری
	if !verifyLicense() {
		fmt.Println("\n\033[31m##########################################")
		fmt.Println("       LICENSE ERROR: UNAUTHORIZED")
		fmt.Printf("       Your Machine ID: %s\n", getMachineID())
		fmt.Println("   Contact Admin to whitelist your server.")
		fmt.Println("##########################################\033[0m")
		os.Exit(1)
	}
	fmt.Println("✅ License Verified Successfully.")

	// ۲. ادامه کدهای اصلی شما
	mode := flag.String("mode", "", "internal: 'server' or 'client'")
	// ... (ادامه کدهای فلگ و اجرای برنامه)
	flag.Parse()
    
    if *mode != "" {
        // اجرای بخش سرور یا کلاینت
        configureLogging()
        // ...
        return
    }
    showInteractiveMenu()
}

// ... بقیه توابع شما (showInteractiveMenu, runServer, runClient و غیره)
