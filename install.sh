#!/bin/bash
set -e

# تنظیمات مخزن و فایل‌ها
GITHUB_REPO="osafari599-coder/aswwa"
LICENSE_URL="https://raw.githubusercontent.com/osafari599-coder/aswwa/main/allowed_servers.txt"
EXECUTABLE_NAME="phantom"
INSTALL_PATH="/usr/local/bin"
SERVICE_NAME="phantom.service"

print_info() { echo -e "\e[34m[INFO]\e[0m $1"; }
print_success() { echo -e "\e[32m[SUCCESS]\e[0m $1"; }
print_error() { echo -e "\e[31m[ERROR]\e[0m $1" >&2; exit 1; }

clear
print_info "Verifying Server Authorization..."
MACHINE_ID=$(hostname)

# ۱. چک کردن لایسنس آنلاین
ALLOWED_LIST=$(curl -sSL "$LICENSE_URL")
if ! echo "$ALLOWED_LIST" | grep -qxw "$MACHINE_ID"; then
    echo -e "\e[31m❌ ACCESS DENIED! Machine ID: $MACHINE_ID \e[0m"
    exit 1
fi
print_success "Access Granted for $MACHINE_ID"

# ۲. دانلود آخرین نسخه (Binary)
print_info "Downloading Phantom binary..."
# در اینجا فرض بر این است که فایل باینری را در Releases گذاشته‌ای
# اگر فایل باینری نداری، باید سورس را دانلود و بیلد کنی
curl -sSLf -o "$EXECUTABLE_NAME" "https://github.com/${GITHUB_REPO}/raw/main/phantom" 
chmod +x "$EXECUTABLE_NAME"
mv "$EXECUTABLE_NAME" "$INSTALL_PATH/"

# ۳. تنظیمات اولیه (رفع ارور Too few arguments)
print_info "Configuring database..."
PANEL_PORT="8080"
PANEL_USER="admin"
PANEL_PASS="admin"

# اجرای دستور ستاپ با آرگومان‌های کامل
$INSTALL_PATH/$EXECUTABLE_NAME --setup-port="$PANEL_PORT" --setup-user="$PANEL_USER" --setup-pass="$PANEL_PASS"

# ۴. ایجاد فایل سرویس
print_info "Creating systemd service..."
cat > "/etc/systemd/system/${SERVICE_NAME}" <<EOF
[Unit]
Description=Phantom Tunnel Service
After=network.target

[Service]
ExecStart=${INSTALL_PATH}/${EXECUTABLE_NAME} --start-panel
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now ${SERVICE_NAME}

# ۵. نمایش خروجی نهایی
if systemctl is-active --quiet $SERVICE_NAME; then
    IP=$(curl -s https://ifconfig.me)
    echo -e "\n\e[32m============================================================\e[0m"
    echo -e "   🚀 PHANTOM TUNNEL IS INSTALLED AND RUNNING!"
    echo -e "============================================================\e[0m"
    echo -e "🔗 Panel URL:  \e[36mhttp://$IP:$PANEL_PORT\e[0m"
    echo -e "👤 Username:   \e[33m$PANEL_USER\e[0m"
    echo -e "🔑 Password:   \e[33m$PANEL_PASS\e[0m"
    echo -e "\e[32m============================================================\e[0m"
else
    print_error "Service failed to start."
fi
