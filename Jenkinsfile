pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        sh 'docker build -t zibby/tplink_r450_prometheus .'
      }
    }
  }
}