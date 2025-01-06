#!groovy
@Library(['github.com/cloudogu/ces-build-lib@1.62.0', 'github.com/cloudogu/zalenium-build-lib@v2.1.1'])
import com.cloudogu.ces.cesbuildlib.*

// Creating necessary git objects, object cannot be named 'git' as this conflicts with the method named 'git' from the library
gitWrapper = new Git(this, "cesmarvin")
gitWrapper.committerName = 'cesmarvin'
gitWrapper.committerEmail = 'cesmarvin@cloudogu.com'
gitflow = new GitFlow(this, gitWrapper)
github = new GitHub(this, gitWrapper)
changelog = new Changelog(this)

// Configuration of repository
repositoryOwner = "cloudogu"
repositoryName = "cesapp-lib"
project = "github.com/${repositoryOwner}/${repositoryName}"
goVersion = "1.23.4"

// Configuration of branches
productionReleaseBranch = "main"
developmentBranch = "develop"
currentBranch = "${env.BRANCH_NAME}"

node('docker') {
    timestamps {
        properties([
                // Keep only the last x builds to preserve space
                buildDiscarder(logRotator(numToKeepStr: '10')),
                // Don't run concurrent builds for a branch, because they use the same workspace directory
                disableConcurrentBuilds(),
        ])

        stage('Checkout') {
            checkout scm
            make 'clean'
        }

        withBuildDependencies {
            stage('Build') {
                make 'compile'
            }

            stage('Unit Tests') {
                make 'unit-test'
                junit allowEmptyResults: true, testResults: 'target/unit-tests/*-tests.xml'
            }

            stage('Integration Test') {
                // If SKIP_DOCKER_TESTS is true, tests which need Docker containers are skipped
                make 'integration-test'
                junit allowEmptyResults: true, testResults: 'target/integration-tests/*-tests.xml'
            }

            stage("Review dog analysis") {
                stageStaticAnalysisReviewDog()
            }
        }

        stage('SonarQube') {
            stageStaticAnalysisSonarQube()
        }

        stageAutomaticRelease()
    }
}

void stageStaticAnalysisReviewDog() {
    def commitSha = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

    withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'sonarqube-gh', usernameVariable: 'USERNAME', passwordVariable: 'REVIEWDOG_GITHUB_API_TOKEN']]) {
        withEnv(["CI_PULL_REQUEST=${env.CHANGE_ID}", "CI_COMMIT=${commitSha}", "CI_REPO_OWNER=${repositoryOwner}", "CI_REPO_NAME=${repositoryName}"]) {
            make 'static-analysis'
        }
    }
}

void stageStaticAnalysisSonarQube() {
    def scannerHome = tool name: 'sonar-scanner', type: 'hudson.plugins.sonar.SonarRunnerInstallation'
    withSonarQubeEnv {
        gitWrapper.fetch()

        if (currentBranch == productionReleaseBranch) {
            echo "This branch has been detected as the production branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (currentBranch == developmentBranch) {
            echo "This branch has been detected as the development branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else if (env.CHANGE_TARGET) {
            echo "This branch has been detected as a pull request."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.pullrequest.key=${env.CHANGE_ID} -Dsonar.pullrequest.branch=${env.CHANGE_BRANCH} -Dsonar.pullrequest.base=${developmentBranch}"
        } else if (currentBranch.startsWith("feature/")) {
            echo "This branch has been detected as a feature branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME}"
        } else {
            echo "This branch has been detected as a miscellaneous branch."
            sh "${scannerHome}/bin/sonar-scanner -Dsonar.branch.name=${env.BRANCH_NAME} "
        }
    }
    timeout(time: 2, unit: 'MINUTES') { // Needed when there is no webhook for example
        def qGate = waitForQualityGate()
        if (qGate.status != 'OK') {
            unstable("Pipeline unstable due to SonarQube quality gate failure")
        }
    }
}

void stageAutomaticRelease() {
    if (!gitflow.isReleaseBranch()) {
        return
    }

    potentiallyCreateDoguDocPR()

    String releaseVersion = gitWrapper.getSimpleBranchName()

    stage('Finish Release') {
        gitflow.finishRelease(releaseVersion, productionReleaseBranch)
    }

    stage('Add Github-Release') {
        releaseId = github.createReleaseWithChangelog(releaseVersion, changelog, productionReleaseBranch)
    }
}

void make(String makeArgs) {
    sh "make ${makeArgs}"
}

void withBuildDependencies(Closure closure) {
    def etcdImage = docker.image('quay.io/coreos/etcd:v3.2.32')
    def etcdContainerName = "${JOB_BASE_NAME}-${BUILD_NUMBER}".replaceAll("\\/|%2[fF]", "-")
    withDockerNetwork { buildnetwork ->
        etcdImage.withRun("--network ${buildnetwork} --name ${etcdContainerName}", 'etcd --listen-client-urls http://0.0.0.0:4001 --advertise-client-urls http://0.0.0.0:4001')
                {
                    new Docker(this)
                            .image("golang:${goVersion}")
                            .mountJenkinsUser()
                            .inside("--network ${buildnetwork} -e ETCD=${etcdContainerName} -e SKIP_SYSLOG_TESTS=true -e SKIP_DOCKER_TESTS=true --volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}")
                                    {
                                        closure.call()
                                    }
                }
    }
}

void potentiallyCreateDoguDocPR() {
    def targetDoguDocDir = "target/dogu-doc"
    def coreDoguChapter = "compendium_en.md"
    def oldCoreDoguChapter = "${targetDoguDocDir}/docs/core/${coreDoguChapter}"
    def newCoreDoguChapter = "target/compendium_en.md"
    def doguDocRepoParts = "cloudogu/dogu-development-docs"
    def doguDocRepo = "github.com/${doguDocRepoParts}.git"
    def doguDocTargetBranch = "main"
    def newBranchName = "feature/update_compendium_after_release_${currentBranch}"
    String releaseVersion = gitWrapper.getSimpleBranchName()
    def gomarkVersion = "v0.4.1-8"

    new Docker(this)
            .image("golang:${goVersion}") // gomarkdoc needs /go/doc/comment from go 1.19+
            .mountJenkinsUser()
            .inside("--volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}") {

                stage('Pull old dogu doc page') {
                    checkout changelog: false, poll: false, scm: scmGit(branches: [[name: doguDocTargetBranch]], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: targetDoguDocDir]], gitTool: 'Default', userRemoteConfigs: [[credentialsId: 'cesmarvin', url: "https://${doguDocRepo}"]])
                }

                stage('Build new dogu doc page') {
                    sh "go install github.com/cloudogu/gomarkdoc/cmd/gomarkdoc@${gomarkVersion}"
                    sh "gomarkdoc --output ${newCoreDoguChapter} core/dogu_v2.go --include-files dogu_v2.go,dogu_v2_security.go,dogu_v2_configuration_field.go,dogu_v2_dependency.go,dogu_v2_versions.go,dogu_v2_volume.go"
                }

                def shouldCreateDoguDocsPR = false

                stage('Compare dogu docs') {
                    // ignore stderr output here, diffing non-existing files always leads to a line count of zero
                    def diffResult = sh(returnStdout: true, script: "diff ${newCoreDoguChapter} ${oldCoreDoguChapter} | wc -l").toString().trim()
                    if ((diffResult as Integer) > 0) {
                        shouldCreateDoguDocsPR = true
                    }
                }

                if (shouldCreateDoguDocsPR) {
                    stage('Create dogu docs PR') {
                        def ghIssueTitle = ":memo: Update and translate dogu.json compendium"
                        def ghIssueBody = "The cesapp-lib release ${releaseVersion} has changed the core.Dogu documentation. Please review the change and translate it. :blue_heart:"
                        def ghPRTitle ="Update to core.Dogu compendium"
                        def ghPRBody1 ="This PR resolves"
                        def ghPRBody2 ="from cesapp-lib version ${releaseVersion}"
                        def ghPRCommentBody = "- [ ] This PR includes a translation :speech_balloon:"

                        sh "cd ${targetDoguDocDir} && git checkout -b ${newBranchName}"
                        // create potential diff by overwriting the original file
                        sh "cp ${newCoreDoguChapter} ${oldCoreDoguChapter}"

                        sh "cd ${targetDoguDocDir} && git add docs/core/${coreDoguChapter}"

                        withCredentials([usernamePassword(credentialsId: 'cesmarvin', usernameVariable: 'GIT_AUTH_USR', passwordVariable: 'GIT_AUTH_PSW')]) {
                            def issueResult = sh(returnStdout: true, script: """curl -L \
                                --write-out '%{http_code}' \
                                -X POST \
                                -H "Accept: application/vnd.github+json" \
                                -u "${GIT_AUTH_USR}:${GIT_AUTH_PSW}" \
                                -H "X-GitHub-Api-Version: 2022-11-28" \
                                -d '{"title":"${ghIssueTitle}","body":"${ghIssueBody}","labels":["enhancement"]}' \
                                https://api.github.com/repos/${doguDocRepoParts}/issues""")

                            // avoid adding more containers (yq is missing here) and mange the issue id by shell foo
                            def issueID = sh(returnStdout: true, script:  """echo -n '${issueResult}' | grep -m 1 '"number":'""")
                            issueID=issueID.replaceAll(/.+:\s+(\d+),/, "\$1").replace("\n", "")
                            sh "echo 'Found issue ->#${issueID}<-'"

                            sh "cd ${targetDoguDocDir} && " +
                                    "git config user.email ces-marvin@cloudogu.com && " +
                                    "git config user.name ces-marvin && " +
                                    "git commit -m '#${issueID} Update to core.Dogu compendium'"
                            sh "cd ${targetDoguDocDir} && git remote set-url origin https://{GIT_AUTH_USR}:${GIT_AUTH_PSW}@${doguDocRepo}"

                            sh "cd ${targetDoguDocDir} && git push --set-upstream origin ${newBranchName}"

                            def pullRequestResult=sh(returnStdout: true, script: """curl -L \
                              --write-out '%{http_code}' \
                              -X POST \
                              -H "Accept: application/vnd.github+json" \
                              -u "${GIT_AUTH_USR}:${GIT_AUTH_PSW}" \
                              -H "X-GitHub-Api-Version: 2022-11-28" \
                              -d '{"title":"${ghPRTitle}","body":"${ghPRBody1} #${issueID} ${ghPRBody2}","head":"${newBranchName}","base":"${doguDocTargetBranch}"}' \
                              https://api.github.com/repos/${doguDocRepoParts}/pulls""")
                            def pullRequestID = sh(returnStdout: true, script:  """echo -n '${pullRequestResult}' | grep -m 1 '"number":'""")
                            pullRequestID=pullRequestID.replaceAll(/.+:\s+(\d+),/, "\$1").replace("\n", "")
                            sh "echo 'Found pullRequestID ->#${pullRequestID}<-'"

                            sh """curl -L \
                              --write-out '%{http_code}' \
                              -X POST \
                              -H "Accept: application/vnd.github+json" \
                              -u "${GIT_AUTH_USR}:${GIT_AUTH_PSW}" \
                              -H "X-GitHub-Api-Version: 2022-11-28" \
                              -d '{"body":"${ghPRCommentBody}"}' \
                              https://api.github.com/repos/${doguDocRepoParts}/issues/${pullRequestID}/comments"""
                        }
                    }
                }
            }
}
