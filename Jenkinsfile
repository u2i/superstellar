discardOldBuilds()

stage('Checkout') {
    node('superstellar-docker-17.12') {
        withCleanup {
            checkout scm
            stash 'source'
        }
    }
}

stage('Build & Test') {
    parallel(
        backend: {
            node('superstellar-docker-17.12') {
                withCleanup {
                    unstash 'source'

                    sh "docker build -t superstellar-backend-builder:${env.BUILD_NUMBER} --target builder -f docker/backend/Dockerfile ."
                    sh "docker build -t superstellar-backend-test:${env.BUILD_NUMBER} --build-arg VERSION=${env.BUILD_NUMBER} -f docker/backend/Dockerfile.test_image ."
                    sh "docker run --rm superstellar-backend-test:${env.BUILD_NUMBER}"

                    masterBranchOnly {
                        stage('Build & Publish backend') {
                            sh "docker build -t superstellar-backend:${env.BUILD_NUMBER} -f docker/backend/Dockerfile ."
                            sh "docker tag superstellar-backend:${env.BUILD_NUMBER} gcr.io/kubernetes-playground-195112/superstellar-backend:${env.BUILD_NUMBER}"

                            withDockerLoggedIntoGCR {
                                sh "docker push gcr.io/kubernetes-playground-195112/superstellar-backend:${env.BUILD_NUMBER}"
                            }
                        }
                    }
                }
            }
        },
        frontend: {
            node('superstellar-docker-17.12') {
                withCleanup {
                    unstash 'source'

                    docker.build("superstellar-frontend:${env.BUILD_NUMBER}", "-f docker/frontend/Dockerfile.production .")

                    masterBranchOnly {
                        stage('Publish frontend') {
                            sh "docker tag superstellar-frontend:${env.BUILD_NUMBER} gcr.io/kubernetes-playground-195112/superstellar-frontend:${env.BUILD_NUMBER}"

                            withDockerLoggedIntoGCR {
                                sh "docker push gcr.io/kubernetes-playground-195112/superstellar-frontend:${env.BUILD_NUMBER}"
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

        node('superstellar-docker-17.12') {
            withCleanup {
                unstash 'source'

                sh "docker build -t superstellar-deployment:${env.BUILD_NUMBER} -f docker/deployment/Dockerfile ."
                withCredentials([file(credentialsId: '5bc94dd2-0a14-4bba-bfd9-f628512b3158', variable: 'FILE')]) {
                    sh 'cp $FILE deployment_volume/service_account.json'
                    sh "docker run -v ${pwd()}/deployment_volume:/deployment_volume superstellar-deployment:${env.BUILD_NUMBER} /deployment_volume/script.sh ${env.BUILD_NUMBER}"
                }
            }
        }
    }
}

stage('Cleanup') {
    node('superstellar-docker-17.12') {
        withCleanup {
            sh 'yes | docker system prune -a --volumes'
        }
    }
}

def withDockerLoggedIntoGCR(Closure cl) {
    withCredentials([file(credentialsId: '5bc94dd2-0a14-4bba-bfd9-f628512b3158', variable: 'FILE')]) {
        sh 'cat $FILE | docker login -u _json_key --password-stdin https://gcr.io'
    }
    cl()
}

def masterBranchOnly(Closure cl) {
    if (env.BRANCH_NAME == 'master') {
        cl()
    }
}
