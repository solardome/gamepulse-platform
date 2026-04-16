FROM golang:1.26 AS build

WORKDIR /src

RUN apt-get update && apt-get install -y --no-install-recommends \
    make \
    protobuf-compiler \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1

COPY . .

RUN make generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/accounts-web ./accounts-web/cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/accounts-web /accounts-web

ENTRYPOINT ["/accounts-web"]
