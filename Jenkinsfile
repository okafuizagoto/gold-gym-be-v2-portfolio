pipeline {
    agent any

    environment {
        PROJECT_ID = "gold-gym-project"

        NAME = "gold-gym-be-v2"
        ORG = "gold-gym"
        K8S_CLUSTER = "gold-gym"

        ARTIFACT_REGISTRY = "localhost:5000"
        DOCKER_REGISTRY = "localhost:5000"
        DOCKER_REGISTRY_URL = "http://${ARTIFACT_REGISTRY}"
        DOCKER_REGISTRY_PROJECT_URL = "${ARTIFACT_REGISTRY}/${PROJECT_ID}"

        PIPELINE_NAME = "Gold Gym Backend"
        PIPELINE_BOT_EMAIL = "goldgym.bot@gmail.com"
        PIPELINE_BOT_NAME = "Gold Gym Pipeline Bot"

        DISCORD_WEBHOOK_URL = "https://discordapp.com/api/webhooks/CHANGE_ME/CHANGE_ME"
    }

    options {
        skipDefaultCheckout(true)
    }

    stages {

        stage('Checkout SCM') {
            steps {
                deleteDir()
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: '*/main']],
                    userRemoteConfigs: [[
                        url: 'https://github.com/okafuizagoto/gold-gym-be-v2.git'
                    ]]
                ])
            }
        }

        stage('Setup Discord Notifications') {
            steps {
                echo 'Discord webhook ready'
            }
        }

        stage('Compile') {
            steps {
                sh 'go version'
                sh 'make build'
            }
        }

        stage('Version') {
            steps {
                script {
                    env.VERSION = sh(
                        script: "git describe --tags --always || echo ${env.BUILD_NUMBER}",
                        returnStdout: true
                    ).trim()
                    echo "VERSION: ${env.VERSION}"
                }
            }
        }

        stage('Dockerize') {
            steps {
                script {
                    echo '> Building Docker image'
                    sh "docker build -t ${NAME}:${env.VERSION} ."
                    sh "docker tag ${NAME}:${env.VERSION} ${NAME}:latest"

                    echo '> Loading image to kind cluster'
                    sh "kind load docker-image ${NAME}:${env.VERSION} --name gold-gym || true"
                    sh "kind load docker-image ${NAME}:latest --name gold-gym || true"
                }
            }
        }

        stage('Vault Chart Injector') {
            steps {
                echo 'Vault disabled (local)'
            }
        }

        stage('Helm Charts') {
            steps {
                script {
                    sh """
                    export KUBECONFIG=/var/jenkins_home/.kube_config

                    echo '> Update Helm values'
                    sed -i 's|repository: draft|repository: ${NAME}|g' charts/values.yaml
                    sed -i 's/tag: dev/tag: ${env.VERSION}/g' charts/values.yaml
                    sed -i 's/name: NAME/name: ${NAME}/g' charts/Chart.yaml

                    echo '> Deploying via Helm'
                    helm upgrade --install ${NAME} charts \
                      --namespace staging \
                      --create-namespace \
                      --set vault.enabled=false \
                      --wait --timeout 240s
                    """
                }
            }
        }

        stage('Helm Charts CHC') {
            steps {
                echo 'CHC deployment skipped (not used)'
            }
        }
    }

    post {

        always {
            cleanWs()
        }

        success {
            sh """
            curl -s -X POST "${env.DISCORD_WEBHOOK_URL}" \
              -H "Content-Type: application/json" \
              -d '{
                "embeds": [{
                  "title": "${env.PIPELINE_NAME} #${env.BUILD_NUMBER} SUCCESS",
                  "color": 65280,
                  "fields": [
                    {"name": "Version", "value": "${env.VERSION}", "inline": true},
                    {"name": "Image", "value": "${NAME}:${env.VERSION}", "inline": true}
                  ]
                }]
              }'
            """
        }

        failure {
            sh """
            curl -s -X POST "${env.DISCORD_WEBHOOK_URL}" \
              -H "Content-Type: application/json" \
              -d '{
                "embeds": [{
                  "title": "${env.PIPELINE_NAME} #${env.BUILD_NUMBER} FAILED",
                  "color": 16711680,
                  "fields": [
                    {"name": "Version", "value": "${env.VERSION}", "inline": true}
                  ]
                }]
              }'
            """
        }
    }
}