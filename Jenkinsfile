def COLOR_MAP = [
    'SUCCESS': 'good',
    'FAILURE': 'danger',
]

def cause = currentBuild.getBuildCauses('hudson.model.Cause$UserIdCause')

pipeline {
    agent {
        kubernetes {
            label 'bexs-slave'
            defaultContainer 'bexs-slave'
        }
    }
    environment {
        //Variaveis utilizadas pelo Slack notification
        // variavel teste: 0=Sucesso, 1=falha
        doError = '0'
        //--------------------------------------------
        SHORT_COMMIT = "${GIT_COMMIT[0..7]}"
        BRANCH_DEV = "develop"
        BRANCH_PROD = "master"
    }    
    stages {
            stage ("Docker Build"){
                steps{
                    container('bexs-slave'){
                        sh "docker build -t gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT} -t temis:${SHORT_COMMIT} ."
                    }
                }
            }
            stage ("Docker Push"){
                steps{
                    container('bexs-slave'){
                        sh "docker push gcr.io/bexs-platform/temis-compliance:${SHORT_COMMIT}"
                    }
                }
            }
            stage ("Delivery for Development"){
                when {
                    branch env.BRANCH_DEV
                }
                steps{
                    container('bexs-slave'){
                        sh 'cd scripts && chmod +x deliver-for-development.sh && ./deliver-for-development.sh'
                    }
                }
            }
            stage ("Functional Tests"){
                  when {
                      branch env.BRANCH_DEV
                  }
                  steps{
                      build job: "/Digital/temis-functional-tests/master", wait: true
                  }
            }
            stage ("Deliver for Sandbox"){
                when {
                    branch env.BRANCH_PROD
                }
                steps{
                    container('bexs-slave'){
                        sh 'cd scripts && chmod +x deliver-for-sandbox.sh && ./deliver-for-sandbox.sh'
                    }
                }
            }
            stage ("System Tests"){
                  when {
                      branch env.BRANCH_PROD
                  }
                  steps{
                      build job: "/Digital/temis-system-tests/master", wait: true
                  }
            }
            stage ("Deliver for Production"){
                when {
                    branch env.BRANCH_PROD
                }
                steps{
                    container('bexs-slave'){
                        sh 'cd scripts &&  chmod +x deliver-for-production.sh && ./deliver-for-production.sh'
                    }
                }
            }
            stage ("Smoke Tests"){
                  steps{
                      sh "sh smoke_test.sh ${BRANCH_NAME}"
                  }
            }
    }
    post {
        always {
            script {
                echo "userName: ${cause.userName}"
            }
            slackSend channel: '#digital-pipelines',
                color: COLOR_MAP[currentBuild.currentResult],
                message: "*${currentBuild.currentResult}:* Job ${env.JOB_NAME} build ${env.BUILD_NUMBER} by ${cause.userName}\n Mais informacoes: ${env.BUILD_URL}"
        }
    }
}
