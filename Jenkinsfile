pipeline {
  environment {
    registry = "zibby/re450-exporter"
    registryCredential = 'f8a79f84-5ad0-43e4-b32c-87e2c6001a62'
    dockerImage = ''
  }
  agent any
  stages {
    stage('Clone Git') {
      steps {
        git 'https://github.com/Zibby/re450-prometheus-exporter'
      }
    }
    stage('Build Image') {
      steps {
        script {
          dockerImage = docker.build registry + ":" + "$env.BRANCH_NAME"
        }
      }
    }
    stage('Push Image') {
      steps {
        script {
          docker.withRegistry( '', registryCredential) {
            dockerImage.push()
          }
          if ("$env.BRANCH_NAME"=='master') {
            docker.withRegistry( '', registryCredential) {
              dockerImage.push('latest')
            }
          }
        }
      }
    }
    stage('clean_up') {
      steps {
        sh "docker rmi $registry:$env.BRANCH_NAME"
      }
    }
  }
  post {
    success {
      slackSend(botUser: true, color: '#36a64f', message: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
    }
    failure {
      slackSend(botUser: true, color: '#b70000', message: "FAIL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})")
    }
  }
}