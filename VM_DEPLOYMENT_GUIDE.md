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
sudo apt install -y git curl build-essential tmux
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

## Step 6: Configure and Install Frontend Dependencies

```bash
# Navigate to frontend directory
cd ~/ExpanceTracker/frontend

# Install dependencies
npm install

# Create .env file for development server
cat > .env << EOF
NODE_ENV=development
VUE_APP_API_URL=http://localhost:1234/
EOF
```

---

## Step 7: Run Backend in Tmux Session

```bash
# Navigate to project root
cd ~/ExpanceTracker

# Create a new tmux session for backend
tmux new-session -d -s backend

# Run the backend in the tmux session
tmux send-keys -t backend './expensetracker' C-m

# Verify backend is running
tmux list-sessions

# (Optional) Attach to see backend logs
# tmux attach -t backend
# Press Ctrl+B then D to detach
```

---

## Step 8: Run Frontend in Tmux Session

```bash
# Navigate to frontend directory
cd ~/ExpanceTracker/frontend

# Create a new tmux session for frontend
tmux new-session -d -s frontend

# Run the frontend dev server in the tmux session
tmux send-keys -t frontend 'npm run serve' C-m

# Verify frontend is running
tmux list-sessions

# (Optional) Attach to see frontend logs
# tmux attach -t frontend
# Press Ctrl+B then D to detach
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

# List all tmux sessions
tmux list-sessions

# Attach to backend session to see logs
tmux attach -t backend
# Press Ctrl+B then D to detach

# Attach to frontend session to see logs
tmux attach -t frontend
# Press Ctrl+B then D to detach
```

---

## Accessing the Application

Once deployed, you can access the application:

- **Frontend**: `http://91.98.145.193:8080`
- **Backend API**: `http://91.98.145.193:1234`

---

## Useful Management Commands

### Tmux Session Management
```bash
# List all tmux sessions
tmux list-sessions

# Attach to backend session
tmux attach -t backend

# Attach to frontend session
tmux attach -t frontend

# Detach from session (while inside tmux)
# Press: Ctrl+B then D

# Kill backend session
tmux kill-session -t backend

# Kill frontend session
tmux kill-session -t frontend

# Kill all sessions
tmux kill-server
```

### Restart Services
```bash
# Restart backend
tmux kill-session -t backend
cd ~/ExpanceTracker
tmux new-session -d -s backend
tmux send-keys -t backend './expensetracker' C-m

# Restart frontend
tmux kill-session -t frontend
cd ~/ExpanceTracker/frontend
tmux new-session -d -s frontend
tmux send-keys -t frontend 'npm run serve' C-m
```

### Update Application
```bash
# Navigate to project directory
cd ~/ExpanceTracker

# Pull latest changes
git pull

# Rebuild backend
go build -o expensetracker main.go

# Update frontend dependencies if needed
cd frontend
npm install
cd ..

# Restart backend
tmux kill-session -t backend
tmux new-session -d -s backend
tmux send-keys -t backend './expensetracker' C-m

# Restart frontend
tmux kill-session -t frontend
cd frontend
tmux new-session -d -s frontend
tmux send-keys -t frontend 'npm run serve' C-m
```

---

## Troubleshooting

### Backend not starting
```bash
# Check if backend tmux session exists
tmux list-sessions | grep backend

# Attach to backend session to see errors
tmux attach -t backend

# Check if port 1234 is in use
sudo netstat -tulpn | grep 1234

# Manually run backend to see errors
cd ~/ExpanceTracker
./expensetracker
```

### Frontend not accessible
```bash
# Check if frontend tmux session exists
tmux list-sessions | grep frontend

# Attach to frontend session to see errors
tmux attach -t frontend

# Check if port 8080 is in use
sudo netstat -tulpn | grep 8080

# Manually run frontend to see errors
cd ~/ExpanceTracker/frontend
npm run serve
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

1. **Enable firewall**: Ensure UFW or iptables is properly configured
2. **Regular updates**: Keep system and dependencies updated
3. **Backup database**: Regularly backup `expenses.db`
4. **Use environment variables**: Store sensitive configuration in env files
5. **Restrict SSH**: Use key-based authentication and disable password login
6. **Monitor tmux sessions**: Regularly check that both services are running
7. **Use screen or tmux with auto-restart**: Consider adding auto-restart logic

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
- Default frontend port is 8080 (configured in Vue.js dev server)
- Both services run in separate tmux sessions for easy management
- Tmux sessions persist even after you logout (unless you kill them or reboot)
- To see logs, attach to the respective tmux session

## Support

For issues or questions, check the logs by attaching to tmux sessions:

```bash
# Backend logs
tmux attach -t backend
# Press Ctrl+B then D to detach

# Frontend logs
tmux attach -t frontend
# Press Ctrl+B then D to detach
```

## Auto-Start on Reboot (Optional)

To automatically start both services after a system reboot, add to crontab:

```bash
# Edit crontab
crontab -e

# Add these lines:
@reboot sleep 10 && cd /home/$(whoami)/ExpanceTracker && tmux new-session -d -s backend './expensetracker'
@reboot sleep 15 && cd /home/$(whoami)/ExpanceTracker/frontend && tmux new-session -d -s frontend 'npm run serve'
```
