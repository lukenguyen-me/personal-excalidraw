# Deployment Guide

Simple guide to deploy Personal Excalidraw to your own server.

## What You Need

- A server/VPS with Ubuntu (or similar Linux)
- At least 1GB RAM and 10GB disk space
- Docker installed (we'll show you how)

## Step 1: Install Docker

If Docker is not installed on your server:

```bash
# Update system
sudo apt update

# Install Docker
curl -fsSL https://get.docker.com -o install-docker.sh
sudo sh install-docker.sh

# Allow your user to run Docker
sudo usermod -aG docker $USER

# Log out and back in for changes to take effect
```

Verify Docker is installed:

```bash
docker --version
```

The script automatically detects and works with both `docker compose` (newer) and `docker-compose` (older) commands.

## Step 2: Get the Application

```bash
# Clone the repository
git clone https://github.com/lukenugyen-me/personal-excalidraw.git
cd personal-excalidraw
```

## Step 3: Configure

### Option A: Interactive Setup (Recommended - No need to edit files!)

Run the interactive configuration:

```bash
./deploy.sh config
```

This will guide you through all the settings step-by-step. Just press Enter to accept defaults or auto-generate secure passwords!

**That's it!** The script will create your configuration file automatically.

### Option B: Manual Setup (For advanced users)

If you prefer to edit the file manually:

```bash
cp .env.production.example .env.production
nano .env.production
```

**You MUST change these two values:**

1. `DB_PASSWORD` - Your database password
2. `ACCESS_KEY` - Your access key to login

To generate secure random values:

```bash
# This will give you two random strings
openssl rand -base64 32
openssl rand -base64 32
```

Copy the first one to `DB_PASSWORD` and the second to `ACCESS_KEY`.

Save the file (Ctrl+X, then Y, then Enter).

## Step 4: Deploy

Start the application:

```bash
./deploy.sh start
```

That's it! Your application is now running at `http://your-server-ip:8080`

## Common Commands

```bash
./deploy.sh config     # Interactive setup (first time)
./deploy.sh start      # Start the application
./deploy.sh stop       # Stop the application
./deploy.sh restart    # Restart the application
./deploy.sh status     # Check if it's running
./deploy.sh logs       # See what's happening
./deploy.sh logs -f    # Watch logs in real-time
./deploy.sh upgrade    # Update to latest version
./deploy.sh backup     # Backup your drawings
./deploy.sh help       # Show all commands
```

## Daily Usage

### Check if Everything is Running

```bash
./deploy.sh status
```

You should see both `app` and `postgres` as "Up" and "healthy".

### View Logs (if something isn't working)

```bash
./deploy.sh logs -f
```

Press Ctrl+C to stop watching.

### Update to Latest Version

```bash
./deploy.sh upgrade
```

This will:

1. Pull the latest code
2. Backup your database automatically
3. Rebuild and restart

### Backup Your Data

```bash
./deploy.sh backup
```

Backups are saved in the `backups/` folder with timestamps.

### Restore from Backup

```bash
# List available backups
ls backups/

# Restore a specific backup
./deploy.sh restore backups/backup_20240115_120000.sql
```

## Setting Up HTTPS (Optional but Recommended)

For a secure connection with a domain name:

### 1. Point Your Domain to Your Server

In your domain registrar (GoDaddy, Namecheap, etc.), create an A record pointing to your server's IP address.

### 2. Install Nginx

```bash
sudo apt install nginx certbot python3-certbot-nginx -y
```

### 3. Create Nginx Configuration

```bash
sudo nano /etc/nginx/sites-available/excalidraw
```

Paste this (replace `yourdomain.com` with your actual domain):

```nginx
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Save and enable it:

```bash
sudo ln -s /etc/nginx/sites-available/excalidraw /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 4. Get SSL Certificate (Free)

```bash
sudo certbot --nginx -d yourdomain.com
```

Follow the prompts. Certbot will automatically configure HTTPS for you.

Now visit `https://yourdomain.com` - you should see your application with a secure connection!

### 5. Update CORS (Important!)

Edit your `.env.production`:

```bash
nano .env.production
```

Change:

```bash
CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

Restart:

```bash
./deploy.sh restart
```

## Troubleshooting

### "Port already in use"

Another service is using port 8080. Either:

- Stop the other service, or
- Change `APP_PORT` in `.env.production` to a different port (like 8081)

### "Cannot connect to database"

Check the logs:

```bash
./deploy.sh logs postgres
```

Make sure the postgres container is running:

```bash
./deploy.sh status
```

### "Access denied" when logging in

Make sure the `ACCESS_KEY` in your `.env.production` matches what you're typing in the login screen.

### Application won't start

View the logs to see the error:

```bash
./deploy.sh logs
```

Common fixes:

1. Make sure Docker is running: `docker ps`
2. Check you have enough disk space: `df -h`
3. Verify your `.env.production` file exists and has the required values

### Still having issues?

```bash
# Stop everything
./deploy.sh stop

# Start fresh
./deploy.sh start

# Watch what happens
./deploy.sh logs -f
```

## Automatic Backups (Recommended)

Set up a daily backup using cron:

```bash
# Edit crontab
crontab -e

# Add this line (backs up daily at 2 AM)
0 2 * * * cd /path/to/personal-excalidraw && ./deploy.sh backup
```

Replace `/path/to/personal-excalidraw` with the actual path.

## Need More Help?

- **Quick Commands**: See [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
- **Technical Details**: See [ARCHITECTURE.md](ARCHITECTURE.md)
- **All Commands**: Run `./deploy.sh help`

## Summary

That's it! The key commands you'll use are:

- `./deploy.sh config` - Set up configuration (first time only)
- `./deploy.sh start` - Start the app
- `./deploy.sh status` - Check if it's running
- `./deploy.sh upgrade` - Update to latest version
- `./deploy.sh backup` - Backup your data

Everything else is automated:

- ✅ Configuration wizard guides you through setup
- ✅ Database created automatically
- ✅ Migrations run automatically
- ✅ Secure passwords generated automatically
- ✅ Application starts automatically

No need to edit configuration files manually unless you want to!

Enjoy your personal drawing application! 🎨
