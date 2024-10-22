# 定义变量
BINARY_NAME=jepsenFuzz
SRC_DIR=cmd/jepsenFuzz
BIN_DIR=bin
GO_FILES=$(wildcard $(SRC_DIR)/*.go)

# 默认目标
all: build

# 检查并创建 bin 目录
check-bin-dir:
	@if [ ! -d $(BIN_DIR) ]; then \
		echo "Creating bin directory..."; \
		mkdir -p $(BIN_DIR); \
	fi
	
# 编译目标
build: $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/jepsenFuzz


# 清理目标
clean:
	@echo "Cleaning up..."
	rm -f $(BIN_DIR)/$(BINARY_NAME)

# 运行目标
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BIN_DIR)/$(BINARY_NAME)

.PHONY: all build clean run
