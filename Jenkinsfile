pipeline {
    agent any
    
    tools {
        nodejs 'NodeJS'
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', 
                    url: 'https://github.com/aqilashraaf/wira-ranking-dashboard-new.git',
                    credentialsId: 'github-pat'
            }
        }
        
        stage('Prepare Environment') {
            steps {
                sh '''
                    # Clean workspace with sudo
                    sudo rm -rf frontend/node_modules frontend/dist backend/wira-backend
                    
                    # Set ownership of workspace to jenkins
                    sudo chown -R jenkins:jenkins .
                    
                    # Set ownership of npm cache to jenkins
                    sudo chown -R jenkins:jenkins /var/lib/jenkins/.npm
                    sudo chmod -R 755 /var/lib/jenkins/.npm
                    
                    # Configure npm
                    sudo -u jenkins npm config set registry https://registry.npmmirror.com/
                    sudo -u jenkins npm config set strict-ssl false
                    sudo -u jenkins npm cache clean --force
                '''
            }
        }
        
        stage('Install Dependencies') {
            steps {
                dir('frontend') {
                    sh '''
                        # Clean as root
                        sudo rm -rf node_modules package-lock.json
                        
                        # Install as jenkins
                        sudo -u jenkins npm install --verbose --no-audit --no-fund --legacy-peer-deps
                    '''
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh '''
                        # Build as jenkins user
                        sudo -u jenkins npm run build
                    '''
                }
            }
        }
        
        stage('Build Backend') {
            steps {
                dir('backend') {
                    sh '''
                        export GO111MODULE=on
                        go mod tidy
                        CGO_ENABLED=0 GOOS=linux go build -o wira-backend
                    '''
                }
            }
        }
        
        stage('Deploy') {
            steps {
                sh '''
                    # Deploy frontend
                    sudo rm -rf /var/www/wira-dashboard/*
                    sudo cp -r frontend/dist/* /var/www/wira-dashboard/
                    sudo chown -R www-data:www-data /var/www/wira-dashboard
                    sudo chmod -R 755 /var/www/wira-dashboard
                    
                    # Deploy backend
                    sudo systemctl stop wira-backend || true
                    sudo mkdir -p /opt/wira-backend
                    sudo cp backend/wira-backend /opt/wira-backend/
                    sudo cp backend/.env /opt/wira-backend/ || true
                    sudo chown -R jenkins:jenkins /opt/wira-backend
                    sudo chmod -R 755 /opt/wira-backend
                    
                    # Deploy Nginx config
                    sudo cp nginx/nginx.conf /etc/nginx/sites-available/default
                    sudo nginx -t && sudo systemctl restart nginx
                    
                    # Start backend service
                    sudo systemctl start wira-backend
                '''
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}
