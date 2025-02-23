package main

import (
	"bytes"
	"context"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streadway/amqp"
)

var (
	amqpURL       = os.Getenv("AMQP_URL")
	queueName     = "image_tasks"
	minioEndpoint = os.Getenv("MINIO_ENDPOINT")
	minioAccess   = os.Getenv("MINIO_ACCESS_KEY")
	minioSecret   = os.Getenv("MINIO_SECRET_KEY")
	minioBucket   = "processed-images"
	maxRetries    = 5
)

func connectRabbitMQ() *amqp.Connection {
	var conn *amqp.Connection
	var err error
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			log.Println("RabbitMQに接続成功")
			return conn
		}
		log.Printf("RabbitMQ接続失敗: %s (再試行 %d/%d)", err, i+1, maxRetries)
		time.Sleep(5 * time.Second)
	}
	log.Fatalf("RabbitMQに接続できませんでした: %s", err)
	return nil
}
func main() {
	conn := connectRabbitMQ()
	defer conn.Close()
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}
	if minioEndpoint == "" {
		minioEndpoint = "minio:9000"
	}
	if minioAccess == "" {
		minioAccess = "minioadmin"
	}
	if minioSecret == "" {
		minioSecret = "minioadmin"
	}

	// RabbitMQ接続
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("RabbitMQ接続失敗: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("チャネルオープン失敗: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("キュー宣言失敗: %s", err)
	}

	msgs, err := ch.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("メッセージ受信登録失敗: %s", err)
	}

	// MinIOクライアント初期化
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccess, minioSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO接続失敗: %s", err)
	}

	ctx := context.Background()
	// バケット作成（存在しない場合）
	if err = minioClient.MakeBucket(ctx, minioBucket, minio.MakeBucketOptions{}); err != nil {
		// 既に存在する場合はエラーを無視
		log.Printf("Bucket作成: %v", err)
	}

	log.Println("ワーカー起動、メッセージ待機中...")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println("画像処理開始")
			// 画像デコード
			srcImg, err := imaging.Decode(bytes.NewReader(d.Body))
			if err != nil {
				log.Printf("画像デコード失敗: %s", err)
				continue
			}

			// リサイズ（幅800、縦はアスペクト比維持）
			resized := imaging.Resize(srcImg, 800, 0, imaging.Lanczos)
			// サムネイル生成（200x200）
			thumbnail := imaging.Thumbnail(srcImg, 200, 200, imaging.Lanczos)

			// JPEGエンコード
			var bufResized, bufThumb bytes.Buffer
			if err = jpeg.Encode(&bufResized, resized, nil); err != nil {
				log.Printf("リサイズ画像エンコード失敗: %s", err)
				continue
			}
			if err = jpeg.Encode(&bufThumb, thumbnail, nil); err != nil {
				log.Printf("サムネイルエンコード失敗: %s", err)
				continue
			}

			// MinIOへアップロード（ここでは固定のファイル名例。実際はユニークな名前にする等の工夫が必要）
			_, err = minioClient.PutObject(ctx, minioBucket, "resized.jpg", &bufResized, int64(bufResized.Len()), minio.PutObjectOptions{ContentType: "image/jpeg"})
			if err != nil {
				log.Printf("リサイズ画像アップロード失敗: %s", err)
				continue
			}
			_, err = minioClient.PutObject(ctx, minioBucket, "thumbnail.jpg", &bufThumb, int64(bufThumb.Len()), minio.PutObjectOptions{ContentType: "image/jpeg"})
			if err != nil {
				log.Printf("サムネイルアップロード失敗: %s", err)
				continue
			}

			log.Println("画像処理完了、MinIOにアップロード")
		}
	}()

	<-forever
}
