# Set these to the desired values
ARTIFACT_ID=cesapp-lib
VERSION=0.3.0

GOTAG?=1.18.1
MAKEFILES_VERSION=6.1.0
LINT_VERSION=v1.45.2
GO_BUILD_FLAGS?=-mod=vendor -a ./...

include build/make/variables.mk
INTEGRATION_TEST_NAME_PATTERN=.*_inttest$$

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/release.mk