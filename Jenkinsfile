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
                git branch: 'main',
                    url: 'https://github.com/aqilashraaf/wira-ranking-dashboard-new.git',
                    credentialsId: 'github-pat'
            }
        }
        
        stage('Build Frontend') {
            steps {
                dir('frontend') {
                    sh 'npm install'
                    sh 'npm run build'
                }
            }
        }
        
        stage('Build Backend') {
            steps {
                dir('backend') {
                    sh 'go build -o wira-backend'
                }
            }
        }

        stage('Build and Push Docker Images') {
            steps {
                script {
                    // Login to private registry
                    sh 'echo $DOCKER_CREDENTIALS_PSW | docker login 173.212.239.58:5000 -u $DOCKER_CREDENTIALS_USR --password-stdin'
                    
                    // Build and push frontend
                    dir('frontend') {
                        sh '''
                            docker build -t 173.212.239.58:5000/frontend:latest .
                            docker push 173.212.239.58:5000/frontend:latest
                        '''
                    }
                    
                    // Build and push backend
                    dir('backend') {
                        sh '''
                            docker build -t 173.212.239.58:5000/backend:latest .
                            docker push 173.212.239.58:5000/backend:latest
                        '''
                    }
                }
            }
        }
        
        stage('Deploy to VPS') {
            steps {
                script {
                    withCredentials([sshUserPrivateKey(credentialsId: 'vps-ssh-key', keyFileVariable: 'SSH_KEY')]) {
                        sh '''
                            # Set up SSH key
                            mkdir -p ~/.ssh
                            cp "$SSH_KEY" ~/.ssh/id_rsa
                            chmod 600 ~/.ssh/id_rsa
                            
                            # Create necessary directories on VPS
                            ssh -o StrictHostKeyChecking=no root@173.212.239.58 'mkdir -p /root/wira-dashboard'
                            ssh -o StrictHostKeyChecking=no root@173.212.239.58 'mkdir -p /root/wira-dashboard/backend'
                            
                            # Copy files to VPS
                            scp -o StrictHostKeyChecking=no docker-compose.yml root@173.212.239.58:/root/wira-dashboard/
                            scp -o StrictHostKeyChecking=no backend/.env.production root@173.212.239.58:/root/wira-dashboard/backend/.env
                            
                            # Deploy on VPS
                            ssh -o StrictHostKeyChecking=no root@173.212.239.58 "cd /root/wira-dashboard && \
                                docker login 173.212.239.58:5000 && \
                                docker-compose pull && \
                                docker-compose down --remove-orphans && \
                                docker-compose up -d && \
                                docker system prune -f"
                                
                            # Clean up SSH key
                            rm -f ~/.ssh/id_rsa
                        '''
                    }
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
    }
}
