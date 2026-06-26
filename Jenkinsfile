pipeline {
    agent any
    stages{
        stage('code quality and testing'){
            step{
                dir('internal/shortener'){
                sh 'go test -v'
            }
            }
        }
        stage('build'){
            step{
                dir('deployments'){
                sh 'docker compose build '
            }
            }   
        }
        stage('deploy'){
            step{
                 dir('deployments'){
                sh 'docker compose up -d'
            }
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished.'
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