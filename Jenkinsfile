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
                    # Clean workspace
                    rm -rf frontend/node_modules frontend/dist backend/wira-backend
                    
                    # Configure npm registries with fallbacks
                    npm config set registry https://registry.npmjs.org/
                    npm config set @vue:registry https://registry.npmjs.org/
                    npm config set @vitejs:registry https://registry.npmjs.org/
                    npm config set strict-ssl false
                    npm cache clean --force
                '''
            }
        }
        
        stage('Install Dependencies') {
            steps {
                dir('frontend') {
                    sh '''
                        # Clean install
                        rm -rf node_modules package-lock.json
                        
                        # Install dependencies
                        npm install --no-audit --no-fund --legacy-peer-deps
                    '''
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh 'npm run build'
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
                    rm -rf /var/www/wira-dashboard/*
                    cp -r frontend/dist/* /var/www/wira-dashboard/
                    chown -R www-data:www-data /var/www/wira-dashboard
                    chmod -R 755 /var/www/wira-dashboard
                    
                    # Deploy backend
                    systemctl stop wira-backend || true
                    mkdir -p /opt/wira-backend
                    cp backend/wira-backend /opt/wira-backend/
                    cp backend/.env /opt/wira-backend/ || true
                    chown -R jenkins:jenkins /opt/wira-backend
                    chmod -R 755 /opt/wira-backend
                    
                    # Deploy Nginx config
                    cp nginx/nginx.conf /etc/nginx/sites-available/default
                    nginx -t && systemctl restart nginx
                    
                    # Start backend service
                    systemctl start wira-backend
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
