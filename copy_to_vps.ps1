# Convert the script to Unix line endings
$content = Get-Content -Path "vps_setup.sh" -Raw
$content = $content -replace "`r`n", "`n"
Set-Content -Path "vps_setup.sh" -Value $content -NoNewline

# Use SSH to copy the file and execute it
$script = @"
scp vps_setup.sh root@173.212.239.58:/root/
ssh root@173.212.239.58 'chmod +x /root/vps_setup.sh && ./vps_setup.sh'
"@

# Save the script to a batch file
Set-Content -Path "run_deploy.bat" -Value $script

# Execute the batch file
Start-Process "cmd.exe" -ArgumentList "/c run_deploy.bat" -NoNewWindow -Wait
