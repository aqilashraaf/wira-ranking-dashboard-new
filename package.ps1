# Create a temporary directory for deployment files
New-Item -ItemType Directory -Force -Path deploy_package
Set-Location deploy_package

# Copy necessary files
Copy-Item ..\docker-compose.yml .
Copy-Item -Recurse ..\nginx .
Copy-Item ..\deploy.sh .

# Create zip archive
Compress-Archive -Path * -DestinationPath ..\deploy.zip -Force

# Clean up
Set-Location ..
Remove-Item -Recurse -Force deploy_package
