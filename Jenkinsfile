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
                    
                    # Configure npm with multiple registries and timeouts
                    npm config set registry https://registry.npmmirror.com/
                    npm config set fetch-retry-mintimeout 100000
                    npm config set fetch-retry-maxtimeout 600000
                    npm config set fetch-timeout 600000
                    npm config set strict-ssl false
                    
                    # Set npm mirrors for specific scopes
                    npm config set @vue:registry https://registry.npmmirror.com/
                    npm config set @vitejs:registry https://registry.npmmirror.com/
                    
                    # Clear npm cache
                    npm cache clean --force
                    
                    # Create .npmrc in the workspace
                    echo "registry=https://registry.npmmirror.com/" > .npmrc
                    echo "strict-ssl=false" >> .npmrc
                    echo "fetch-retry-mintimeout=100000" >> .npmrc
                    echo "fetch-retry-maxtimeout=600000" >> .npmrc
                    echo "fetch-timeout=600000" >> .npmrc
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
                        
                        # Install dependencies with increased network timeout
                        npm install --no-audit --no-fund --legacy-peer-deps --network-timeout 600000
                    '''
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh 'NODE_OPTIONS="--max-old-space-size=4096" npm run build'
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
