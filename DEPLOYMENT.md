# ExpenseTracker Server Configuration

## 1. Frontend (Static Files)
Serve the `frontend/dist/` directory on port 8080

### Using Python:
```bash
cd frontend/dist
python3 -m http.server 8080
```

### Using Node.js serve:
```bash
npm install -g serve
serve -s frontend/dist -l 8080
```

### Using nginx:
```nginx
server {
    listen 8080;
    server_name your-domain.com;
    
    location / {
        root /path/to/ExpenseTracker/frontend/dist;
        try_files $uri $uri/ /index.html;
    }
}
```

## 2. Backend (Go API)
Run the Go binary on port 1234

```bash
./expensetracker
```

## 3. Environment Variables
Update the frontend/.env file with your server IP:
```
NODE_ENV=production
VUE_APP_API_URL=http://YOUR_SERVER_IP:1234/
```

## 4. Firewall Rules
Make sure these ports are open:
- Port 8080 (Frontend)
- Port 1234 (Backend API)

## 5. Database
The SQLite database will be created automatically in the same directory as the binary.
