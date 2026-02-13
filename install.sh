#!/bin/bash

set -e

# تنظیمات
GITHUB_REPO="webwizards-team/Phantom-Tunnel"
LICENSE_URL="https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"
EXECUTABLE_NAME="phantom"
INSTALL_PATH="/usr/local/bin"
SERVICE_NAME="phantom.service"
WORKING_DIR="/etc/phantom"

# توابع چاپ
print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_success() { echo -e "\e[32m[SUCCESS]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }

clear
print_info "Starting Phantom Tunnel Installation..."

if [ "$(id -u)" -ne 0 ]; then
  print_error "This script must be run as root."
fi

# ۱. چک کردن لایسنس آنلاین
print_info "Checking Server Authorization..."
MACHINE_ID=$(hostname)
ALLOWED_LIST=$(curl -sSL "$LICENSE_URL")

if ! echo "$ALLOWED_LIST" | grep -qxw "$MACHINE_ID"; then
    echo -e "\e[31m--------------------------------------------\e[0m"
    echo -e "❌ ACCESS DENIED!"
    echo -e "Your Machine ID: \e[32m$MACHINE_ID\e[0m"
    echo -e "This server is not in the allowed list."
    echo -e "Please send your Machine ID to Admin."
    echo -e "\e[31m--------------------------------------------\e[0m"
    exit 1
fi

print_success "Server Authorized: $MACHINE_ID"

# ۲. نصب پیش‌نیازها
print_info "Installing dependencies..."
apt-get update -y > /dev/null && apt-get install -y curl grep > /dev/null

# ۳. دانلود و نصب (بقیه کدهای قبلی خودت)
ARCH=$(uname -m)
[ "$ARCH" == "x86_64" ] && ASSET_NAME="phantom-amd64" || ASSET_NAME="phantom-arm64"

LATEST_TAG=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep -oP '"tag_name": "\K[^"]+')
DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_TAG}/${ASSET_NAME}"

print_info "Downloading binary..."
curl -sSLf -o "$EXECUTABLE_NAME" "$DOWNLOAD_URL"
chmod +x "$EXECUTABLE_NAME"
mkdir -p "$WORKING_DIR"
mv "$EXECUTABLE_NAME" "$INSTALL_PATH/"

# ۴. ساخت سرویس
cat > "/etc/systemd/system/${SERVICE_NAME}" <<EOF
[Unit]
Description=Phantom Tunnel Service
After=network.target

[Service]
ExecStart=${INSTALL_PATH}/${EXECUTABLE_NAME} --start-panel
WorkingDirectory=${WORKING_DIR}
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now ${SERVICE_NAME}

print_success "Phantom Tunnel installed and started!"
