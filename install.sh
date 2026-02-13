#!/bin/bash
set -e

LICENSE_URL="https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"

print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }

clear
print_info "Verifying Server..."
MACHINE_ID=$(hostname)

# چک کردن آنلاین
if ! curl -sSL "$LICENSE_URL" | grep -qxw "$MACHINE_ID"; then
    echo -e "\e[31m"
    echo "❌ ACCESS DENIED!"
    echo "Your Machine ID: $MACHINE_ID"
    echo "Please send this ID to Admin for access."
    echo -e "\e[0m"
    exit 1
fi

print_info "Access Granted. Starting Phantom Tunnel Installation..."
# ادامه کدهای دانلود و نصب خودت...
# بخش انتهای install.sh برای نمایش اطلاعات ورود
if systemctl is-active --quiet $SERVICE_NAME; then
    IP=$(curl -s https://ifconfig.me)
    # پورت را از تنظیمات یا مقدار پیش‌فرض ۸۰۸۰ بردار
    PORT=${PANEL_PORT:-8080}
    
    echo -e "\n\e[32m============================================================\e[0m"
    echo -e "   🚀 PHANTOM TUNNEL IS INSTALLED AND RUNNING!"
    echo -e "============================================================\e[0m"
    echo -e "🔗 Panel URL:  \e[36mhttp://$IP:$PORT\e[0m"
    echo -e "👤 Username:   \e[33m${PANEL_USER:-admin}\e[0m"
    echo -e "🔑 Password:   \e[33m${PANEL_PASS:-admin}\e[0m"
    echo -e "\e[32m============================================================\e[0m"
else
    print_error "Service failed to start. Check logs with: journalctl -u phantom -f"
fi
