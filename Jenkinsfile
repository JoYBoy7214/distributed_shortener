pipeline {
    agent any
    
    stages {
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
            sh 'docker system prune -f'
        }
        success {
            echo 'DEPLOYMENT SUCCESSFUL: URL Shortener is live!'
        }
        failure {
            echo 'DEPLOYMENT FAILED: Check logs. Stack was not updated.'
        }
    }
}