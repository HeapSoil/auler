
# ==============================================================================
# 默认执行 all 目标
.DEFAULT_GOAL := all

# make 时默认执行以下伪目标
.PHONY: all
all: go.format  go.build

# ==============================================================================

include scripts/make-rules/common.mk # 在 common.mk 中定义全局 Makefile 变量方便后面引用
include scripts/make-rules/golang.mk
include scripts/make-rules/code.mk
include scripts/make-rules/tools.mk



define USAGE_OPTIONS

Options:
  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="miniblog test"
  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.
  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS


# ==============================================================================

# Makefile 代码生成

.PHONY: ca
ca: ## 生成 CA 文件.
	@$(MAKE) gen.ca

.PHONY: protoc
protoc: ## 编译 protobuf 文件.
	@$(MAKE) gen.protoc

.PHONY: deps
deps: ## 安装依赖，例如：生成需要的代码、安装需要的工具等.
	@$(MAKE) gen.deps


# ==============================================================================

# Makefile build 二进制

.PHONY: build
build: go.tidy  ## 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@$(MAKE) go.build


# ==============================================================================

# Makefile 清理构建产物

.PHONY: clean
clean: ## 清理构建产物、临时文件等.
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)


# ==============================================================================

# Makefile Hack/Tools

.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy


.PHONY: format
format:  ## 格式化 Go 源码.
	@$(MAKE) go.format


.PHONY: help
help: Makefile ## 打印 Makefile help 信息.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<TARGETS> <OPTIONS>\033[0m\n\n\033[35mTargets:\033[0m\n"} /^[0-9A-Za-z._-]+:.*?##/ { printf "  \033[36m%-45s\033[0m %s\n", $$1, $$2 } /^\$$\([0-9A-Za-z_-]+\):.*?##/ { gsub("_","-", $$1); printf "  \033[36m%-45s\033[0m %s\n", tolower(substr($$1, 3, length($$1)-7)), $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' Makefile #$(MAKEFILE_LIST)
	@echo -e "$$USAGE_OPTIONS"