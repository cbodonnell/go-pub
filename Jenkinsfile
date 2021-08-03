pipeline {
    agent any
    environment {
        GOROOT = "${tool type: 'go', name: 'go1.15.6'}/go"
    }
    stages {
        stage('build') {
            steps {
                echo 'building...'
                sh 'echo $GOROOT'
                sh '$GOROOT/bin/go build'
            }
        }
        stage('test') {
            steps {
                echo 'testing...'
            }
        }
        stage('deploy') {
            steps {
                echo 'deploying...'
                sh 'sudo systemctl stop go-pub'
                sh 'sudo cp go-pub /etc/go-pub/go-pub'
                sh 'sudo systemctl start go-pub'
            }
        }
    }
    post {
        cleanup {
            deleteDir()
        }
    }
}