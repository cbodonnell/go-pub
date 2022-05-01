pipeline {
    agent any
    stages {
        stage('build') {
            steps {
                echo 'building...'
                sh 'sudo chmod +x ./scripts/build && ./scripts/build'
            }
        }
        stage('test') {
            steps {
                echo 'testing...'
                sh 'npm run test:prod'
            }
        }
        stage('deploy') {
            steps {
                echo 'deploying...'
                sh 'sudo chmod +x ./scripts/deploy && ./scripts/deploy'
            }
        }
    }
    post {
        cleanup {
            deleteDir()
        }
    }
}