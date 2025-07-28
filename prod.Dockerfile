FROM golang:1.24-bookworm AS builder

COPY ./go.mod ./go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o /go/bin/portfolio-v2025-api .

FROM scratch

# バイナリを配置
COPY --from=builder /go/bin/portfolio-v2025-api /go-bin/portfolio-v2025-api

# 非特権ユーザで実行
USER 10000

EXPOSE 80

CMD ["/go-bin/portfolio-v2025-api"]
