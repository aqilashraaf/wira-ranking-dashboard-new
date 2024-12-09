#!/bin/bash

# Create a temporary directory for deployment files
mkdir -p deploy_package
cd deploy_package

# Copy necessary files
cp ../docker-compose.yml .
cp -r ../nginx .
cp ../deploy.sh .

# Create tar archive
tar -czf ../deploy.tar.gz *

# Clean up
cd ..
rm -rf deploy_package
