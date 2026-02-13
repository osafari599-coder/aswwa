#!/bin/bash

set -e

# --- ۱. تنظیمات متغیرها ---
GITHUB_REPO="webwizards-team/Phantom-Tunnel"
EXECUTABLE_NAME="phantom"
INSTALL_PATH="/usr/local/bin"
SERVICE_NAME="phantom.service"
WORKING_DIR="/etc/phantom"

# --- ۲. تعریف توابع (حتماً باید قبل از استفاده تعریف شوند) ---
print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_success() { echo -e "\e[32m[SUCCESS]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }
print_warning() { echo -e "\e[33m⚠️ WARNING: $1\033[0m"; }

# --- ۳. شروع نصب ---
clear
print_info "Starting Phantom Tunnel Installation..."

if [ "$(id -u)" -ne 0 ]; then
  print_error "This script must be run as root. Please use 'sudo'."
fi

# ایجاد پوشه کاری
mkdir -p "$WORKING_DIR"

# --- ۴. بخش لایسنسینگ ---
print_info "Checking License..."
MACHINE_ID=$(hostname)

if [ ! -f "$WORKING_DIR/license.key" ]; then
    echo -e "\e[33m--------------------------------------------\e[0m"
    echo -e "Your Machine ID: \e[32m$MACHINE_ID\e[0m"
    echo -e "Please provide this ID to the provider to get your Key."
    echo -e "\e[33m--------------------------------------------\e[0m"
    
    read -p "Enter your License Key: " USER_KEY
    if [ -z "$USER_KEY" ]; then
        print_error "License Key cannot be empty."
    fi
    echo "$USER_KEY" | sudo tee "$WORKING_DIR/license.key" > /dev/null
    print_success "License key saved."
fi

# --- ۵. بررسی وابستگی‌ها ---
print_info "Checking for dependencies (curl, grep)..."
if command -v apt-get &> /dev/null; then
    apt-get update -y > /dev/null && apt-get install -y -qq curl grep > /dev/null
elif command -v yum &> /dev/null; then
    yum install -y curl grep > /dev/null
fi
print_success "Dependencies are satisfied."

# --- ۶. تشخیص معماری و دانلود فایل ---
ARCH=$(uname -m)
case $ARCH in
    x86_64) ASSET_NAME="phantom-amd64" ;;
    aarch64 | arm64) ASSET_NAME="phantom-arm64" ;;
    *) print_error "Unsupported architecture: $ARCH" ;;
esac

LATEST_TAG=$(curl -s "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep -oP '"tag_name": "\K[^"]+')
if [ -z "$LATEST_TAG" ]; then
    print_error "Failed to fetch the latest release tag."
fi

DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${LATEST_TAG}/${ASSET_NAME}"

print_info "Downloading latest binary..."
TMP_DIR=$(mktemp -d); cd "$TMP_DIR"
if ! curl -sSLf -o "$EXECUTABLE_NAME" "$DOWNLOAD_URL"; then
    print_error "Download failed."
fi

# جایگذاری فایل اجرایی
chmod +x "$EXECUTABLE_NAME"
if systemctl is-active --quiet $SERVICE_NAME; then
    sudo systemctl stop $SERVICE_NAME
fi
mv "$EXECUTABLE_NAME" "$INSTALL_PATH/"
print_success "Binary installed successfully."

# --- ۷. تنظیم سرویس ---
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

print_success "Phantom Tunnel is now RUNNING!"
echo "------------------------------------------------------------"
