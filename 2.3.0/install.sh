#!/bin/bash
set -e

BASE_URL="https://raw.githubusercontent.com/osafari599-coder/aswwa/main"
VERSION="2.3.0"
LICENSE_URL="${BASE_URL}/allowed_servers.txt"
SOURCE_URL="${BASE_URL}/${VERSION}/phantom.go"

print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_success() { echo -e "\e[32m[SUCCESS]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }

clear
print_info "Starting Installation v${VERSION}..."

# ۱. تایید لایسنس
MACHINE_ID=$(hostname)
if ! curl -sSL "$LICENSE_URL" | grep -qxw "$MACHINE_ID"; then
    print_error "ACCESS DENIED! Machine ID: $MACHINE_ID"
fi
print_success "Server Authorized."

# ۲. نصب پیش‌نیازها (Go)
if ! command -v go &> /dev/null; then
    print_info "Installing Go Lang..."
    sudo apt update && sudo apt install golang -y
fi

# ۳. دانلود سورس و بیلد کردن
print_info "Downloading and Compiling..."
curl -sSLf -o "phantom.go" "$SOURCE_URL"
go mod init phantom &> /dev/null || true
go mod tidy &> /dev/null || true
go build -o phantom phantom.go
chmod +x phantom
mv phantom /usr/local/bin/

# ۴. تنظیمات و سرویس
print_info "Finalizing..."
/usr/local/bin/phantom --setup-port=8080 --setup-user=admin --setup-pass=admin

cat > /etc/systemd/system/phantom.service <<EOF
[Unit]
Description=Phantom Service
After=network.target

[Service]
ExecStart=/usr/local/bin/phantom --start-panel
Restart=always

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now phantom.service

print_success "DONE! System is running."
