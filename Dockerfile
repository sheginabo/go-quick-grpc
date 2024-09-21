# 構建階段
FROM --platform=linux/amd64 golang:1.23 AS builder

# 設置工作目錄
WORKDIR /app

# 安裝必要的包和 protobuf 編譯器
RUN apk add --no-cache \
    git \
    go \
    musl-dev \
    protobuf \
    protobuf-dev

# 檢查 protoc 版本
RUN protoc --version

# 安裝 protoc-gen-go 和 protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 設置 PATH 以包含 Go 二進制文件目錄
ENV PATH="$PATH:$(go env GOPATH)/bin"

# 複製 go mod 和 sum 文件 優化緩存
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . ./

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# 最終階段
FROM --platform=linux/amd64 alpine:latest

WORKDIR /app

# 從構建階段複製二進制文件
COPY --from=builder /app/main .

# 安裝 protobuf 運行時依賴
RUN apk add --no-cache protobuf

# 暴露端口
EXPOSE 8080 9090

# 運行
CMD ["/app/main"]