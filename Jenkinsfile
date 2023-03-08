#!groovy
@Library(['github.com/cloudogu/dogu-build-lib@v1.10.0', 'github.com/cloudogu/ces-build-lib@1.57.0'])
import com.cloudogu.ces.cesbuildlib.*
import com.cloudogu.ces.dogubuildlib.*

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
    def etcdImage = docker.image('quay.io/coreos/etcd:v3.2.5')
    def etcdContainerName = "${JOB_BASE_NAME}-${BUILD_NUMBER}".replaceAll("\\/|%2[fF]", "-")
    withDockerNetwork { buildnetwork ->
        etcdImage.withRun("--network ${buildnetwork} --name ${etcdContainerName}", 'etcd --listen-client-urls http://0.0.0.0:4001 --advertise-client-urls http://0.0.0.0:4001')
                {
                    new Docker(this)
                            .image('golang:1.18.6')
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
    def oldCoreDoguChapter="${targetDoguDocDir}/docs/core/${coreDoguChapter}"
    def newCoreDoguChapter = "target/${coreDoguChapter}"
    def doguDocRepo = "https://github.com/cloudogu/dogu-development-docs.git"
    def doguDocTargetBranch ="main"
    // FIXME
    def newBranchName = "platano"

    new Docker(this)
            .image('golang:1.18.6')
            .mountJenkinsUser()
            .inside("--network ${buildnetwork} -e ETCD=${etcdContainerName} -e SKIP_SYSLOG_TESTS=true -e SKIP_DOCKER_TESTS=true --volume ${WORKSPACE}:/go/src/${project} -w /go/src/${project}") {

                stage('Pull old dogu doc page') {
                    checkout changelog: false, poll: false, scm: scmGit(branches: [[name: doguDocTargetBranch]], extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: targetDoguDocDir]], gitTool: 'Default', userRemoteConfigs: [[credentialsId: 'cesmarvin', url: doguDocRepo]])
                }

                stage('Build new dogu doc page') {
                    sh 'gomarkdoc --output ${localCoreDoguChapter} core/dogu_v2.go'
                }

                def shouldCreateDoguDocsPR = false

                stage('Compare dogu docs') {
                    // ignore stderr output here, diffing non-existing files always leads to a line count of zero
                    def diffResult = sh(returnStdout: true, script: "diff ${newCoreDoguChapter} ${oldCoreDoguChapter} | wc -l").toString().trim()
                    if(diffResult > 0) {
                        shouldCreateDoguDocsPR = true
                    }
                }

                if(shouldCreateDoguDocsPR) {
                    stage('Create dogu docs PR') {
                        sh "cd ${targetDoguDocDir} && git checkout -b ${newBranchName}"
                        // create potential diff by overwriting the original file
                        sh "cp ${localCoreDoguChapter} ${oldCoreDoguChapter}"

                        sh "cd ${targetDoguDocDir} && git add ${oldCoreDoguChapter}"

                        sh "cd ${targetDoguDocDir} && git --author='ces-marvin <ces-marvin@cloudogu.com>' commit ${oldCoreDoguChapter}"

                        sh "cd ${targetDoguDocDir} && git push --set-upstream origin feature/${newBranchName}"

                        // create PR from branch -> main "tranl8 plz! (o_0')"
                    }
                }
            }
}