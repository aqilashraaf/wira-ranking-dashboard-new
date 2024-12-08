pipeline {
    agent any
    
    tools {
        nodejs 'NodeJS'
    }
    
    environment {
        NPM_CONFIG_REGISTRY = 'https://registry.npmmirror.com'
        NPM_CONFIG_STRICT_SSL = 'false'
        DOCKER_REGISTRY = '173.212.239.58:5000'
        DOCKER_CREDENTIALS = credentials('docker-credentials')
        SSH_CREDENTIALS = credentials('vps-ssh-key')
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
                    # Update Docker Compose
                    curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)" -o docker-compose
                    chmod +x docker-compose
                    
                    # Clean workspace
                    rm -rf frontend/node_modules frontend/dist backend/wira-backend
                    
                    # Configure npm
                    npm config set registry https://registry.npmmirror.com/
                    npm config set @vue:registry https://registry.npmmirror.com/
                    npm config set @vitejs:registry https://registry.npmmirror.com/
                    npm config set strict-ssl false
                    npm config set fetch-retries 5
                    npm config set fetch-retry-factor 2
                    npm config set fetch-retry-mintimeout 20000
                    npm config set fetch-retry-maxtimeout 120000
                    
                    # Clear npm cache
                    npm cache clean --force
                    
                    # Create .npmrc file
                    echo "registry=https://registry.npmmirror.com/" > .npmrc
                    echo "@vue:registry=https://registry.npmmirror.com/" >> .npmrc
                    echo "@vitejs:registry=https://registry.npmmirror.com/" >> .npmrc
                    echo "strict-ssl=false" >> .npmrc
                    echo "fetch-retries=5" >> .npmrc
                    echo "fetch-retry-factor=2" >> .npmrc
                    echo "fetch-retry-mintimeout=20000" >> .npmrc
                    echo "fetch-retry-maxtimeout=120000" >> .npmrc
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
                        npm install --no-audit --no-fund --legacy-peer-deps --registry https://registry.npmmirror.com/
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
        
        stage('Build Docker Images') {
            steps {
                sh '''
                    # Make sure jenkins can access docker
                    sudo chmod 666 /var/run/docker.sock
                    
                    # Build images
                    ./docker-compose build
                    docker tag wira-ranking-dashboard-new_frontend ${DOCKER_REGISTRY}/frontend:latest
                    docker tag wira-ranking-dashboard-new_backend ${DOCKER_REGISTRY}/backend:latest
                '''
            }
        }
        
        stage('Push Docker Images') {
            steps {
                sh '''
                    echo $DOCKER_CREDENTIALS_PSW | docker login ${DOCKER_REGISTRY} -u $DOCKER_CREDENTIALS_USR --password-stdin
                    docker push ${DOCKER_REGISTRY}/frontend:latest
                    docker push ${DOCKER_REGISTRY}/backend:latest
                '''
            }
        }
        
        stage('Deploy to VPS') {
            steps {
                sshagent(['vps-ssh-key']) {
                    sh '''
                        # Copy docker-compose file to VPS
                        scp -o StrictHostKeyChecking=no docker-compose.yml root@173.212.239.58:/opt/wira-ranking/
                        
                        # Update Docker Compose on VPS and deploy
                        ssh -o StrictHostKeyChecking=no root@173.212.239.58 "
                            cd /opt/wira-ranking
                            curl -L 'https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-$(uname -s)-$(uname -m)' -o docker-compose
                            chmod +x docker-compose
                            ./docker-compose pull
                            ./docker-compose up -d
                        "
                    '''
                }
            }
        }
    }
    
    post {
        always {
            sh 'docker logout ${DOCKER_REGISTRY}'
            cleanWs()
        }
    }
}
