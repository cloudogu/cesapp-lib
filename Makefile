# Set these to the desired values
ARTIFACT_ID=k8s-ces-setup
VERSION=0.0.0

GOTAG?=1.18.1
MAKEFILES_VERSION=5.1.0

include build/make/variables.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/release.mk