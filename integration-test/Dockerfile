# Step 1: Modules caching
FROM golang:1.18.2-alpine3.14 as modules
COPY go.mod go.sum /modules/
COPY go.mod go.sum integration-test/sample.csv /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Tests
FROM golang:1.18.2-alpine3.14
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app

RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=arm64

CMD ["go", "test", "-v", "./integration-test/..."]