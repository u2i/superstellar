import static Constants.*

class Constants {
    static final AWS_CREDENTIALS = 'adae9164-3272-49b1-ab0a-6475983d0ed2'
    static final DOCKER_REGISTRY_URL = 'https://738931564455.dkr.ecr.eu-central-1.amazonaws.com'
    static final DOCKER_REGISTRY_CREDENTIALS = 'u2i-jenkins-ecr-docker-login'
}

discardOldBuilds()

@NonCPS
def jsonParse(def json) {
    new groovy.json.JsonSlurperClassic().parseText(json)
}

stage('Checkout') {
    node {
        withCleanup {
            checkout scm
            stash 'source'
        }
    }
}

stage('Build & Test') {
    parallel(
        backend: {
            node('docker') {
                withCleanup {
                    unstash 'source'

                    docker.build('superstellar-backend-test:latest', '--target builder -f docker/backend/Dockerfile .').inside {
                        sh 'git config --global user.name "Dummy" && git config --global user.email "dummy@example.com"'
                        sh """
                            go get github.com/onsi/ginkgo github.com/onsi/gomega
                            go test superstellar/...
                        """
                    }

                    masterBranchOnly {
                        stage('Package & Publish backend') {
                            def image = docker.build('u2i/superstellar')

                            privateRegistry {
                                image.push(env.BUILD_NUMBER)
                                image.push('latest')
                            }
                        }
                    }
                }
            }
        },
        frontend: {
            node('docker') {
                withCleanup {
                    unstash 'source'

                    def frontend_image = docker.build("superstellar-backend:1.0", "-f docker/frontend/Dockerfile.production .")

                    masterBranchOnly {
                        stage('Package & Publish frontend') {
                            def image = docker.build('u2i/superstellar_nginx')

                            privateRegistry {
                                image.push(env.BUILD_NUMBER)
                                image.push('latest')
                            }
                        }
                    }
                }
            }
        }
    )
}

masterBranchOnly {
    stage(name: 'Deploy') {
        milestone 1

        node('docker') {
            withCleanup {
                String fileName = java.util.UUID.randomUUID().toString()

                aws("--region=eu-central-1 ecs list-tasks --cluster=default > ${fileName}")

                String resultString = readFile(fileName).trim()
                def result = jsonParse(resultString)

                for (String task in result['taskArns']) {
                    String taskId = task.split(":task/")[1]
                    aws("--region=eu-central-1 ecs stop-task --task ${taskId} > /dev/null")
                }
            }
        }
    }

    stage(name: 'Health Check') {
        sleep 10

        retry(5) {
            try {
                node('docker') {
                    withCleanup {
                        sh 'docker run --rm appropriate/curl --fail -I http://superstellar.u2i.is'
                    }
                }
            } catch(Exception e) {
                sleep 10
                throw e
            }
        }
    }
}

def masterBranchOnly(Closure cl) {
    if (env.BRANCH_NAME == 'master') {
        cl()
    }
}

def aws(String cmd) {
    withAwsCredentials(AWS_CREDENTIALS) {
        sh """
            docker run --rm -e AWS_ACCESS_KEY_ID=${env.AWS_ACCESS_KEY_ID} -e AWS_SECRET_ACCESS_KEY=${env.AWS_SECRET_ACCESS_KEY} -v ${pwd()}:/tmp -w /tmp \
                mikesir87/aws-cli:1.11.3 aws $cmd
       """
   }
}

def privateRegistry(Closure cl) {
    String fileName = java.util.UUID.randomUUID().toString()

    aws("--region eu-central-1 ecr get-login > $fileName")

    String dockerLoginCommand = readFile(fileName).trim()

    sh dockerLoginCommand

    docker.withRegistry(DOCKER_REGISTRY_URL) {
        cl()
    }
}
