// ابتدای فایل phantom.go
package main

import (
    "io/ioutil"
    "net/http"
    "os"
    "strings"
    "fmt"
    "time"
    // سایر ایمپورت‌ها...
)

const LicenseURL = "https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"

func verifyLicense() bool {
    client := http.Client{Timeout: 5 * time.Second}
    resp, err := client.Get(LicenseURL)
    if err != nil { return false }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    mID, _ := os.Hostname()
    return strings.Contains(string(body), strings.TrimSpace(mID))
}

func main() {
    if !verifyLicense() {
        fmt.Println("❌ Unauthorized Server!"); os.Exit(1)
    }
    // ادامه کدهای اصلی پنل...
}
