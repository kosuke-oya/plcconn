FROM golang:1.23

WORKDIR /app
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/gopls@latest && \
    go install go.uber.org/mock/mockgen@latest

# COPY go.mod .
# COPY go.sum .

# RUN go mod download

# COPY . .

# RUN go build -o /main .

# CMD ["/main"]