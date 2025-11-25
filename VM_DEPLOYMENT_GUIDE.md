# VM Deployment Guide - Expense Tracker

This guide provides step-by-step commands to deploy the Expense Tracker application on a VM.

## Prerequisites
- Ubuntu/Debian VM with SSH access
- Sudo privileges
- Internet connection

---

## Step 1: Update System and Install Dependencies

```bash
# Update package lists
sudo apt update

# Upgrade existing packages
sudo apt upgrade -y

# Install required packages
sudo apt install -y git curl nginx build-essential
```

---

## Step 2: Install Node.js and npm

```bash
# Download and install Node.js 18.x LTS
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -

# Install Node.js
sudo apt install -y nodejs

# Verify installation
node --version
npm --version
```

---

## Step 3: Install Go

```bash
# Download Go 1.21 (or latest version)
cd /tmp
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# Remove any previous Go installation
sudo rm -rf /usr/local/go

# Extract and install Go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

---

## Step 4: Clone the Repository

```bash
# Navigate to home directory
cd ~

# Clone your repository (replace with your actual repo URL)
git clone https://github.com/DarkBenky/ExpanceTracker.git

# Navigate to project directory
cd ExpanceTracker
```

---

## Step 5: Build the Backend

```bash
# Navigate to project root
cd ~/ExpanceTracker

# Build the Go backend
go build -o expensetracker main.go

# Make it executable
chmod +x expensetracker

# Verify the build
ls -lh expensetracker
```

---

## Step 6: Configure and Build the Frontend

```bash
# Navigate to frontend directory
cd ~/ExpanceTracker/frontend

# Install dependencies
npm install

# Create production .env file
cat > .env.production << EOF
NODE_ENV=production
VUE_APP_API_URL=http://YOUR_VM_IP:1234/
EOF

# Replace YOUR_VM_IP with your actual VM IP address
# Example: VUE_APP_API_URL=http://192.168.1.100:1234/
sed -i 's/YOUR_VM_IP/YOUR_ACTUAL_IP_HERE/g' .env.production

# Build the frontend for production
npm run build

# Verify the build
ls -lh dist/
```

---

## Step 7: Configure Nginx for Frontend

```bash
# Create nginx configuration file
sudo tee /etc/nginx/sites-available/expensetracker << EOF
server {
    listen 8080;
    server_name _;

    root /home/$(whoami)/ExpanceTracker/frontend/dist;
    index index.html;

    location / {
        try_files \$uri \$uri/ /index.html;
    }

    # Enable gzip compression
    gzip on;
    gzip_vary on;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
EOF

# Remove default nginx site
sudo rm -f /etc/nginx/sites-enabled/default

# Enable the expense tracker site
sudo ln -s /etc/nginx/sites-available/expensetracker /etc/nginx/sites-enabled/

# Test nginx configuration
sudo nginx -t

# Restart nginx
sudo systemctl restart nginx

# Enable nginx to start on boot
sudo systemctl enable nginx
```

---

## Step 8: Create Systemd Service for Backend

```bash
# Create systemd service file for the backend
sudo tee /etc/systemd/system/expensetracker.service << EOF
[Unit]
Description=Expense Tracker Backend API
After=network.target

[Service]
Type=simple
User=$(whoami)
WorkingDirectory=/home/$(whoami)/ExpanceTracker
ExecStart=/home/$(whoami)/ExpanceTracker/expensetracker
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd daemon
sudo systemctl daemon-reload

# Start the backend service
sudo systemctl start expensetracker

# Enable the service to start on boot
sudo systemctl enable expensetracker

# Check service status
sudo systemctl status expensetracker
```

---

## Step 9: Configure Firewall

```bash
# Allow SSH (if using UFW)
sudo ufw allow ssh

# Allow HTTP
sudo ufw allow 80/tcp

# Allow frontend port
sudo ufw allow 8080/tcp

# Allow backend API port
sudo ufw allow 1234/tcp

# Enable firewall
sudo ufw --force enable

# Check firewall status
sudo ufw status
```

---

## Step 10: Verify Deployment

```bash
# Check if backend is running
curl http://localhost:1234/

# Check if frontend is accessible
curl http://localhost:8080/

# Check backend service logs
sudo journalctl -u expensetracker -f

# Check nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

---

## Accessing the Application

Once deployed, you can access the application:

- **Frontend**: `http://YOUR_VM_IP:8080`
- **Backend API**: `http://YOUR_VM_IP:1234`

Replace `YOUR_VM_IP` with your VM's actual IP address.

---

## Useful Management Commands

### Backend Service Management
```bash
# Start the backend
sudo systemctl start expensetracker

# Stop the backend
sudo systemctl stop expensetracker

# Restart the backend
sudo systemctl restart expensetracker

# View backend logs
sudo journalctl -u expensetracker -n 100 --no-pager

# Follow backend logs in real-time
sudo journalctl -u expensetracker -f
```

### Nginx Management
```bash
# Start nginx
sudo systemctl start nginx

# Stop nginx
sudo systemctl stop nginx

# Restart nginx
sudo systemctl restart nginx

# Reload nginx (without dropping connections)
sudo systemctl reload nginx

# Test nginx configuration
sudo nginx -t
```

### Update Application
```bash
# Navigate to project directory
cd ~/ExpanceTracker

# Pull latest changes
git pull

# Rebuild backend
go build -o expensetracker main.go

# Rebuild frontend
cd frontend
npm install
npm run build
cd ..

# Restart services
sudo systemctl restart expensetracker
sudo systemctl restart nginx
```

---

## Troubleshooting

### Backend not starting
```bash
# Check service status
sudo systemctl status expensetracker

# View detailed logs
sudo journalctl -u expensetracker -n 50 --no-pager

# Check if port 1234 is in use
sudo netstat -tulpn | grep 1234
```

### Frontend not accessible
```bash
# Check nginx status
sudo systemctl status nginx

# Verify nginx configuration
sudo nginx -t

# Check nginx error logs
sudo tail -n 50 /var/log/nginx/error.log
```

### Database issues
```bash
# Check if database file exists
ls -lh ~/ExpanceTracker/expenses.db

# Ensure correct permissions
chmod 644 ~/ExpanceTracker/expenses.db
```

### Connection issues
```bash
# Check if services are listening on correct ports
sudo netstat -tulpn | grep -E '(1234|8080)'

# Test API connectivity
curl -v http://localhost:1234/

# Check firewall rules
sudo ufw status numbered
```

---

## Security Recommendations

1. **Use HTTPS**: Set up SSL certificates with Let's Encrypt
2. **Change default ports**: Consider using non-standard ports
3. **Enable firewall**: Ensure UFW or iptables is properly configured
4. **Regular updates**: Keep system and dependencies updated
5. **Backup database**: Regularly backup `expenses.db`
6. **Use environment variables**: Store sensitive configuration in env files
7. **Restrict SSH**: Use key-based authentication and disable password login

---

## Optional: Setup HTTPS with Let's Encrypt

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain SSL certificate (replace YOUR_DOMAIN with your actual domain)
sudo certbot --nginx -d YOUR_DOMAIN

# Auto-renewal test
sudo certbot renew --dry-run
```

---

## Backup and Restore

### Backup
```bash
# Create backup directory
mkdir -p ~/backups

# Backup database
cp ~/ExpanceTracker/expenses.db ~/backups/expenses_$(date +%Y%m%d_%H%M%S).db

# Backup entire project
tar -czf ~/backups/expensetracker_$(date +%Y%m%d_%H%M%S).tar.gz ~/ExpanceTracker
```

### Restore
```bash
# Restore database
cp ~/backups/expenses_YYYYMMDD_HHMMSS.db ~/ExpanceTracker/expenses.db

# Restart backend service
sudo systemctl restart expensetracker
```

---

## Notes

- The database file (`expenses.db`) is created automatically on first run
- Default backend port is 1234 (can be changed in the Go code)
- Default frontend port is 8080 (configured in nginx)
- All logs are available via `journalctl` for the backend service
- Nginx logs are in `/var/log/nginx/`

## Support

For issues or questions, check the logs first:
```bash
# Backend logs
sudo journalctl -u expensetracker -f

# Nginx logs
sudo tail -f /var/log/nginx/error.log
```
