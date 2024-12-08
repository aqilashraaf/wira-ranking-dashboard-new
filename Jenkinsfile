pipeline {
    agent any
    
    tools {
        nodejs 'NodeJS'
    }
    
    environment {
        NPM_CONFIG_REGISTRY = 'https://registry.npmmirror.com'
        NPM_CONFIG_STRICT_SSL = 'false'
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
                    
                    # Create .npmrc file
                    cat << EOF > .npmrc
registry=https://registry.npmmirror.com/
@vue:registry=https://registry.npmmirror.com/
@vitejs:registry=https://registry.npmmirror.com/
strict-ssl=false
fetch-retries=5
fetch-retry-factor=2
fetch-retry-mintimeout=20000
fetch-retry-maxtimeout=120000
EOF
                    
                    # Clear npm cache
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
                        
                        # Copy .npmrc to frontend directory
                        cp ../.npmrc .
                        
                        # Install dependencies
                        export NODE_OPTIONS="--max-old-space-size=4096"
                        npm install --prefer-offline --no-audit --no-fund --legacy-peer-deps
                    '''
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh '''
                        export NODE_OPTIONS="--max-old-space-size=4096"
                        npm run build
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
