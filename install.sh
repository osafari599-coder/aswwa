#!/bin/bash
set -e

# آدرس‌دهی بر اساس ساختار مخزن شما
BASE_URL="https://raw.githubusercontent.com/osafari599-coder/aswwa/main"
VERSION="2.3.0"
LICENSE_URL="${BASE_URL}/allowed_servers.txt"
BINARY_URL="${BASE_URL}/${VERSION}/phantom"

print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_success() { echo -e "\e[32m[SUCCESS]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }

clear
print_info "Starting Phantom Tunnel v${VERSION} Installation..."

# ۱. تایید هویت سرور
MACHINE_ID=$(hostname)
ALLOWED_LIST=$(curl -sSL "$LICENSE_URL")

if ! echo "$ALLOWED_LIST" | grep -qxw "$MACHINE_ID"; then
    echo -e "\e[31m❌ ACCESS DENIED! Machine ID: $MACHINE_ID is not authorized.\e[0m"
    exit 1
fi

print_success "Server Authorized: $MACHINE_ID"

# ۲. دانلود فایل اجرایی (Binary)
print_info "Downloading Phantom binary..."
curl -sSLf -o "phantom" "$BINARY_URL" || print_error "Download failed! Ensure 'phantom' binary exists in /${VERSION}/ folder."
chmod +x phantom
mv phantom /usr/local/bin/

# ۳. تنظیمات اولیه
print_info "Configuring service..."
/usr/local/bin/phantom --setup-port=8080 --setup-user=admin --setup-pass=admin

# ۴. نصب سرویس سیستم‌دی
cat > /etc/systemd/system/phantom.service <<EOF
[Unit]
Description=Phantom Tunnel Service
After=network.target

[Service]
ExecStart=/usr/local/bin/phantom --start-panel
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now phantom.service

print_success "Phantom Tunnel is now RUNNING!"
