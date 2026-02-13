#!/bin/bash

set -e

LICENSE_URL="https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"
GITHUB_REPO="webwizards-team/Phantom-Tunnel"
EXECUTABLE_NAME="phantom"
INSTALL_PATH="/usr/local/bin"
SERVICE_NAME="phantom.service"

print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_success() { echo -e "\e[32m[SUCCESS]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }

clear
print_info "Checking Authorization..."
MACHINE_ID=$(hostname)

# بررسی لایسنس آنلاین
ALLOWED_LIST=$(curl -sSL "$LICENSE_URL")
if ! echo "$ALLOWED_LIST" | grep -qxw "$MACHINE_ID"; then
    echo -e "\e[31m"
    echo "❌ ACCESS DENIED!"
    echo "Your Machine ID: $MACHINE_ID"
    echo "This server is not authorized. Please contact the provider."
    echo -e "\e[0m"
    exit 1
fi

print_success "Server Authorized: $MACHINE_ID"

# ادامه مراحل نصب
print_info "Starting installation..."
# ... (کدهای دانلود و نصب سرویس که قبلاً داشتی)
