REPO_ROOT:=${CURDIR}
OUT_DIR=$(REPO_ROOT)/bin

# the output binary name, overridden when cross compiling
KIND_BINARY_NAME?=keink
KIND_BUILD_FLAGS?=-trimpath -ldflags="-buildid= -w -X=sigs.k8s.io/kind/pkg/cmd/kind/version.GitCommit=$(COMMIT)"

export GO111MODULE=on
all: build
# builds keink in a container, outputs to $(OUT_DIR)
keink:
	go build -v -o "$(OUT_DIR)/$(KIND_BINARY_NAME)" $(KIND_BUILD_FLAGS)
# alias for building kind
build: keink


################################################################################
# ============================== Auto-Update ===================================
# update generated code, gofmt, etc.
update:
	hack/make-rules/update/all.sh
# update generated code
generate:
	hack/make-rules/update/generated.sh
# gofmt
gofmt:
	hack/make-rules/update/gofmt.sh
################################################################################
# ================================== Linting ===================================
# run linters, ensure generated code, etc.
verify:
	hack/make-rules/verify/all.sh
# code linters
lint:
	hack/make-rules/verify/lint.sh