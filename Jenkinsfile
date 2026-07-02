pipeline {
    agent any
    
    stages {
        stage('env placing'){
            steps {
                withCredentials([
                    file(credentialsId: 'Distributed_url_shortener_DB_ENV', variable: 'DB_env'),
                    file(credentialsId: 'Distributed_url_shortener_Gateway_env', variable: 'Gateway_env'),
                    file(credentialsId: 'Distributed_url_shortener_shortener_env', variable: 'shortener_env'),
                ]){
                    script{
                        echo 'Placing .env files in required directory'
                        sh "cp \$DB_env ./deployments/DB.env"
                        sh "cp \$Gateway_env ./deployments/gateway.env"
                        sh "cp \$shortener_env ./deployments/shortener.env"
                    }
                }
            }
        }
        stage('Code Quality and Testing') {
            steps {
                // Using ./... at the root level tells Go to find and run tests in ALL folders
                sh 'go test -v ./...'
            }
        }
        
        stage('Build Docker Images') {
            steps {
                dir('deployments') {
                    sh 'docker compose build'
                }
            }
        }
        
        stage('Deploy Services') {
            steps {
                dir('deployments') {
                    // The -d flag is CRITICAL. It runs containers in the background (detached mode)
                    sh 'docker compose up -d'
                    
                }
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished. Cleaning up unused Docker artifacts...'
            // The -f flag forces the prune, bypassing the (y/n) user prompt
            sh 'docker compose down'
            sh 'docker system prune -f'
            sh 'rm -rf ./deployments/DB.env ./deployments/gateway.env ./deployments/shortener.env'
        }
        success {
            echo 'DEPLOYMENT SUCCESSFUL: URL Shortener is live!'
        }
        failure {
            echo 'DEPLOYMENT FAILED: Check logs. Stack was not updated.'
        }
    }
}