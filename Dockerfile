ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    mkdir -p /out && \
    go build -trimpath -o /out/linkedin-mcp .

FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

COPY --from=builder /out/linkedin-mcp /app/linkedin-mcp

EXPOSE 8080

ENTRYPOINT ["/app/linkedin-mcp"]

