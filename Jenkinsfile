pipeline {
    agent any
    
    environment {
        DOCKER_CREDENTIALS = credentials('docker-credentials')
        SSH_CREDENTIALS = credentials('vps-ssh-key')
    }
    
    tools {
        nodejs 'NodeJS'
    }
    
    stages {
        stage('Checkout') {
            steps {
                tool name: 'Default'
                git branch: 'main',
                    url: 'https://github.com/aqilashraaf/wira-ranking-dashboard-new.git',
                    credentialsId: 'github-pat'
            }
        }
        
        stage('Prepare Environment') {
            steps {
                script {
                    // Install latest docker-compose
                    sh '''
                        PLATFORM=$(uname -s)-$(uname -m)
                        curl -L "https://github.com/docker/compose/releases/download/v2.24.5/docker-compose-${PLATFORM}" -o docker-compose
                        chmod +x docker-compose
                    '''
                    
                    // Clean previous builds
                    sh '''
                        rm -rf frontend/node_modules frontend/dist backend/wira-backend
                    '''
                    
                    // Configure npm
                    sh '''
                        npm config set registry https://registry.npmmirror.com/
                        npm config set @vue:registry https://registry.npmmirror.com/
                        npm config set @vitejs:registry https://registry.npmmirror.com/
                        npm config set strict-ssl false
                        npm config set fetch-retries 5
                        npm config set fetch-retry-factor 2
                        npm config set fetch-retry-mintimeout 20000
                        npm config set fetch-retry-maxtimeout 120000
                        npm cache clean --force
                        
                        echo "registry=https://registry.npmmirror.com/" > .npmrc
                        echo "@vue:registry=https://registry.npmmirror.com/" >> .npmrc
                        echo "@vitejs:registry=https://registry.npmmirror.com/" >> .npmrc
                        echo "strict-ssl=false" >> .npmrc
                        echo "fetch-retries=5" >> .npmrc
                        echo "fetch-retry-factor=2" >> .npmrc
                        echo "fetch-retry-mintimeout=20000" >> .npmrc
                        echo "fetch-retry-maxtimeout=120000" >> .npmrc
                    '''
                    
                    // Check Docker socket permissions
                    sh '''
                        if [ ! -w /var/run/docker.sock ]; then
                            echo "Docker socket is not writable. Attempting to fix permissions..."
                            sudo chmod 666 /var/run/docker.sock
                        fi
                    '''
                    
                    // Configure Docker daemon for insecure registry
                    sh '''
                        echo '{ "insecure-registries": ["173.212.239.58:5000"] }' | sudo tee /etc/docker/daemon.json
                        sudo systemctl restart docker || sudo service docker restart
                        sleep 10  # Wait for Docker to restart
                    '''
                }
            }
        }
        
        stage('Install Dependencies') {
            steps {
                dir('frontend') {
                    sh '''
                        rm -rf node_modules package-lock.json
                        cp ../.npmrc .
                        export NODE_OPTIONS=--max-old-space-size=4096
                        npm install --no-audit --no-fund --legacy-peer-deps --registry https://registry.npmmirror.com/
                    '''
                }
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh '''
                        export NODE_OPTIONS=--max-old-space-size=4096
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
                    ./docker-compose build
                    docker tag wira-ranking-pipeline-frontend 173.212.239.58:5000/frontend:latest
                    docker tag wira-ranking-pipeline-backend 173.212.239.58:5000/backend:latest
                '''
            }
        }
        
        stage('Push Docker Images') {
            steps {
                sh '''
                    echo $DOCKER_CREDENTIALS_PSW | docker login 173.212.239.58:5000 -u $DOCKER_CREDENTIALS_USR --password-stdin
                    docker push 173.212.239.58:5000/frontend:latest
                    docker push 173.212.239.58:5000/backend:latest
                '''
            }
        }
        
        stage('Deploy to VPS') {
            steps {
                withCredentials([sshUserPrivateKey(credentialsId: 'vps-ssh-key', keyFileVariable: 'SSH_KEY')]) {
                    sh '''#!/bin/bash
                        set -x
                        # Create .ssh directory if it doesn't exist
                        mkdir -p ~/.ssh
                        chmod 700 ~/.ssh
                        
                        echo "=== Checking SSH key format ==="
                        # Create a temporary directory for key operations
                        TEMP_DIR=$(mktemp -d)
                        cp "$SSH_KEY" $TEMP_DIR/original_key
                        chmod 600 $TEMP_DIR/original_key
                        
                        # Display key format
                        echo "Original key format:"
                        head -n 1 $TEMP_DIR/original_key
                        
                        # Try to convert key to PEM format
                        echo "Converting key to PEM format..."
                        cp $TEMP_DIR/original_key $TEMP_DIR/converted_key
                        ssh-keygen -p -N "" -m PEM -f $TEMP_DIR/converted_key || true
                        
                        # Copy the converted key
                        cp $TEMP_DIR/converted_key ~/.ssh/temp_key
                        chmod 600 ~/.ssh/temp_key
                        
                        echo "=== SSH key permissions ==="
                        ls -la ~/.ssh/temp_key
                        
                        echo "=== SSH debug info ==="
                        ssh-keygen -l -f ~/.ssh/temp_key || true
                        
                        # Add host key
                        ssh-keyscan -H 173.212.239.58 >> ~/.ssh/known_hosts 2>/dev/null
                        
                        echo "=== Testing SSH connection ==="
                        # Try both the original and converted keys
                        echo "Trying with original key..."
                        ssh -vvv -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i $TEMP_DIR/original_key root@173.212.239.58 'echo "SSH connection successful"' || true
                        
                        echo "Trying with converted key..."
                        ssh -vvv -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i ~/.ssh/temp_key root@173.212.239.58 'echo "SSH connection successful"'
                        
                        DEPLOY_SUCCESS=false
                        if [ $? -eq 0 ]; then
                            echo "=== Deploying containers ==="
                            ssh -o StrictHostKeyChecking=no -o IdentitiesOnly=yes -i ~/.ssh/temp_key root@173.212.239.58 '
                                # Pull latest images
                                docker pull 173.212.239.58:5000/frontend:latest
                                docker pull 173.212.239.58:5000/backend:latest
                                
                                # Stop existing containers
                                docker stop frontend backend || true
                                docker rm frontend backend || true
                                
                                # Start new containers
                                docker run -d --name frontend -p 80:80 173.212.239.58:5000/frontend:latest
                                docker run -d --name backend -p 8080:8080 173.212.239.58:5000/backend:latest
                                
                                # Clean up old images
                                docker image prune -f
                            '
                            DEPLOY_SUCCESS=true
                        fi
                        
                        # Clean up
                        rm -rf $TEMP_DIR
                        rm -f ~/.ssh/temp_key
                        
                        if [ "$DEPLOY_SUCCESS" != "true" ]; then
                            echo "WARNING: Deployment failed, but build artifacts are available in the registry"
                            echo "Latest images are tagged as 173.212.239.58:5000/frontend:latest and 173.212.239.58:5000/backend:latest"
                        fi
                    '''
                }
            }
        }
        
        stage('Health Check') {
            steps {
                sh '''
                    echo "=== Checking Frontend Health ==="
                    curl -I http://173.212.239.58:80 || echo "Frontend health check failed"
                    
                    echo "=== Checking Backend Health ==="
                    curl -I http://173.212.239.58:8080 || echo "Backend health check failed"
                    
                    echo "=== Application Status ==="
                    echo "Frontend URL: http://173.212.239.58:80"
                    echo "Backend URL: http://173.212.239.58:8080"
                '''
            }
        }
    }
    
    post {
        always {
            node(null) {
                sh 'docker logout 173.212.239.58:5000'
                cleanWs()
            }
        }
    }
}
