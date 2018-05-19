pipeline {
  agent any

  stages {
    stage('Build') {
      steps {
        sh 'docker build -t thesis .'
      }
    }
    stage('Test') {
      steps {
        sh 'docker run thesis go test ./...'
      }
    }
  }

  post {
    always {
      sh 'docker stop `docker ps -a -q -f status=exited` &> /dev/null || true &> /dev/null'
      sh 'docker rm -v `docker ps -a -q -f status=exited` &> /dev/null || true &> /dev/null'
      sh 'docker rmi `docker images --filter "dangling=true" -q --no-trunc` &> /dev/null || true &> /dev/null'
    }
  }
}