# Personal Excalidraw - Quick Reference

Quick command reference for managing your Personal Excalidraw deployment.

## Deployment Script Commands

```bash
./deploy.sh <command>
```

### Setup & Basic Operations

| Command | Description |
|---------|-------------|
| `config` | Interactive configuration setup (first time) |
| `start` | Start the application |
| `stop` | Stop the application |
| `restart` | Restart the application |
| `status` | Show service status and health |

### Logs & Monitoring

| Command | Description |
|---------|-------------|
| `logs` | Show last 100 lines of logs |
| `logs -f` | Follow logs in real-time |
| `logs app` | Show logs for app service only |
| `logs postgres` | Show logs for postgres service only |
| `logs app -f` | Follow app logs in real-time |

### Maintenance

| Command | Description |
|---------|-------------|
| `upgrade` | Pull latest changes, backup DB, rebuild and restart |
| `backup` | Create database backup |
| `restore <file>` | Restore database from backup file |
| `clean` | Remove all containers and images (keeps database) |

### Help

| Command | Description |
|---------|-------------|
| `help` | Show detailed help and usage |

## Common Workflows

### First Time Setup

```bash
# 1. Interactive configuration (no file editing!)
./deploy.sh config

# 2. Start application
./deploy.sh start

# 3. Check status
./deploy.sh status
```

**Alternative (manual setup):**
```bash
cp .env.production.example .env.production
nano .env.production  # Edit DB_PASSWORD and ACCESS_KEY
./deploy.sh start
```

### Regular Update

```bash
# One command to do everything
./deploy.sh upgrade
```

### View Logs

```bash
# See what's happening
./deploy.sh logs -f
```

### Backup Before Changes

```bash
# Create backup
./deploy.sh backup

# Make changes...

# If needed, restore
./deploy.sh restore backups/backup_20240101_120000.sql
```

### Troubleshooting

```bash
# Check service health
./deploy.sh status

# View recent errors
./deploy.sh logs app

# Restart if needed
./deploy.sh restart
```

## Docker Compose (Alternative)

If you prefer using Docker Compose directly:

```bash
# Start
docker compose -f docker-compose.prod.yml --env-file .env.production up -d

# Stop
docker compose -f docker-compose.prod.yml down

# View logs
docker compose -f docker-compose.prod.yml logs -f

# Rebuild and restart
docker compose -f docker-compose.prod.yml up -d --build
```

## Environment Variables

Key variables in `.env.production`:

```bash
# Required - MUST be changed
DB_PASSWORD=your_strong_password
ACCESS_KEY=your_random_key

# Optional
APP_PORT=8080
CORS_ALLOWED_ORIGINS=https://yourdomain.com
LOG_LEVEL=info
AUTH_ENABLED=true
```

Generate secure values:

```bash
# Database password
openssl rand -base64 32

# Access key
openssl rand -base64 32
```

## Health Check

```bash
# Check if application is running
curl http://localhost:8080/health

# Should return: {"status":"ok"}
```

## File Locations

| Path | Description |
|------|-------------|
| `.env.production` | Production environment configuration |
| `backups/` | Database backup files |
| `docker-compose.prod.yml` | Production Docker Compose configuration |
| `Dockerfile.prod` | Production Docker build file |

## Useful Docker Commands

```bash
# View running containers
docker ps

# View all containers (including stopped)
docker ps -a

# View container resource usage
docker stats

# View container logs
docker logs personal-excalidraw-app -f

# Execute command in container
docker exec -it personal-excalidraw-app sh

# Connect to PostgreSQL
docker exec -it personal-excalidraw-postgres psql -U postgres personal_excalidraw

# Clean up unused Docker resources
docker system prune -f
```

## Port Access

Default ports (can be changed in `.env.production`):

- **Application**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API**: http://localhost:8080/api/*

## Getting Help

- **Deployment Guide**: [DEPLOYMENT.md](DEPLOYMENT.md)
- **Architecture**: [ARCHITECTURE.md](ARCHITECTURE.md)
- **Project Plan**: [PLAN.md](PLAN.md)
- **README**: [README.md](README.md)

For the deployment script help:

```bash
./deploy.sh help
```
