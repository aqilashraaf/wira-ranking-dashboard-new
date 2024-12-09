#!/bin/bash

# Install required packages
apt-get update
apt-get install -y nginx docker.io docker-compose git

# Clone the repository
cd /root
rm -rf wira-ranking-dashboard-new
git clone https://github.com/aqilashraaf/wira-ranking-dashboard-new.git
cd wira-ranking-dashboard-new

# Set up nginx
mkdir -p /etc/nginx/sites-available
mkdir -p /etc/nginx/sites-enabled

# Copy nginx configuration
cp nginx/nginx.conf /etc/nginx/sites-available/wira-dashboard
ln -sf /etc/nginx/sites-available/wira-dashboard /etc/nginx/sites-enabled/

# Remove default nginx site if it exists
rm -f /etc/nginx/sites-enabled/default

# Start the services
docker-compose down || true
docker-compose pull
docker-compose up -d

# Restart nginx
systemctl restart nginx

# Show status
docker-compose ps
systemctl status nginx
