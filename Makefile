include .project/gomod-project.mk
# set DEBUG=1 for debug build, eg. 'DEBUG=1 make build'
ifeq ($(DEBUG),)
BUILD_FLAGS=-ldflags="-s -w"
else
BUILD_FLAGS=-gcflags=all="-N -l"
endif
DEBUG_TAG="$(shell whoami)-test"
# -timeout=3m -v
TEST_FLAGS=-timeout=3m
export GOPRIVATE=github.com/tbilicode
export COVERAGE_EXCLUSIONS="vendor|tests|third_party|api/pb/|main\.go|testsuite\.go|gomock|mocks/|\.gen\.go|\.pb\.go"

.PHONY: *

.SILENT:

default: help

all: clean tools generate version change_log build test

#
# clean produced files
#
clean:
	echo "*** cleaning"
	go clean ./...
	rm -rf \
		${COVPATH} \
		${PROJ_BIN} \
		./mocks

tools:
	echo "*** building tools"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/effective-security/cov-report/cmd/cov-report@latest
	go install go.uber.org/mock/mockgen@latest

version:
	echo "*** building version"
	gofmt -r '"GIT_VERSION" -> "$(GIT_VERSION)"' internal/version/current.template > internal/version/current.go


build:
	echo "*** Building UI client"
	go build ${BUILD_FLAGS} -o ${PROJ_ROOT}/bin/bog ./cmd/bog

change_log:
	echo "Recent changes" > ./change_log.txt
	echo "Build Version: $(GIT_VERSION)" >> ./change_log.txt
	echo "Commit: $(GIT_HASH)" >> ./change_log.txt
	echo "==================================" >> ./change_log.txt
	git log -n 20 --pretty=oneline --abbrev-commit >> ./change_log.txt
