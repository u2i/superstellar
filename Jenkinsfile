node('docker') {
    withCleanup {
        withTimestamps {
            checkout(scm)

            withDockerCompose { compose ->
                compose.createLocalUser('app')

                compose.exec('app', "go get superstellar")
                compose.exec('app', "go test superstellar")
            }
        }
    }
}
