#!/bin/bash

# Install required packages if not present
apt-get update
apt-get install -y nginx docker.io docker-compose unzip

# Create project directory
mkdir -p /root/wira-ranking-dashboard-new
cd /root/wira-ranking-dashboard-new

# Unzip deployment package
unzip -o /root/deploy.zip

# Copy necessary files
cp -r /root/jenkins/workspace/wira-ranking-pipeline/* .

# Create nginx directory and configuration
mkdir -p /etc/nginx/sites-available
mkdir -p /etc/nginx/sites-enabled

# Copy nginx configuration
cp nginx/nginx.conf /etc/nginx/sites-available/wira-dashboard
ln -sf /etc/nginx/sites-available/wira-dashboard /etc/nginx/sites-enabled/

# Remove default nginx site if it exists
rm -f /etc/nginx/sites-enabled/default

# Create frontend directory
mkdir -p /var/www/wira-dashboard

# Start the services
docker-compose down || true
docker-compose pull
docker-compose up -d

# Restart nginx
systemctl restart nginx

# Show status
docker-compose ps
systemctl status nginx
