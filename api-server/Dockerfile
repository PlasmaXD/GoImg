FROM golang:1.20-alpine
WORKDIR /app
# Goモジュールのファイルをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download
# ソースコード全体をコピー
COPY . .
RUN go build -o api-server .
EXPOSE 8080
CMD ["./api-server"]
