package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/streadway/amqp"
)

var (
	amqpURL   string
	queueName = "image_tasks"

	minioEndpoint string
	minioAccess   string
	minioSecret   string
	minioBucket   = "processed-images"

	minioClient *minio.Client
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// POST以外は静的HTML（アップロードフォーム）を返す
		http.ServeFile(w, r, "./static/index.html")
		return
	}

	// フォームから画像ファイルを受け取る
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "フォーム解析エラー", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "ファイル取得エラー", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// ファイルをメモリ上に読み込み
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "ファイル読み込みエラー", http.StatusInternalServerError)
		return
	}

	// RabbitMQへ接続し、ジョブキューに画像データを送信
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Printf("RabbitMQ接続失敗: %s", err)
		http.Error(w, "メッセージブローカー接続エラー", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("チャネルオープン失敗: %s", err)
		http.Error(w, "チャネルオープンエラー", http.StatusInternalServerError)
		return
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
		log.Printf("キュー宣言失敗: %s", err)
		http.Error(w, "キュー宣言エラー", http.StatusInternalServerError)
		return
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        fileBytes,
			Headers: map[string]interface{}{
				"filename": handler.Filename,
			},
		},
	)
	if err != nil {
		log.Printf("メッセージ送信失敗: %s", err)
		http.Error(w, "メッセージ送信エラー", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("画像を受け付けました。処理を開始します。"))
}

func imagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	opts := minio.ListObjectsOptions{Prefix: "", Recursive: true}
	objectCh := minioClient.ListObjects(ctx, minioBucket, opts)

	type imageInfo struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	var images []imageInfo

	for obj := range objectCh {
		if obj.Err != nil {
			log.Printf("ListObjects error: %v", obj.Err)
			continue
		}
		presigned, err := minioClient.PresignedGetObject(ctx, minioBucket, obj.Key, time.Hour, url.Values{})
		if err != nil {
			log.Printf("URL生成失敗: %v", err)
			continue
		}
		images = append(images, imageInfo{Name: obj.Key, URL: presigned.String()})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

func main() {
	// 環境変数からRabbitMQ接続先を取得（未設定ならデフォルト）
	amqpURL = os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	minioEndpoint = os.Getenv("MINIO_ENDPOINT")
	if minioEndpoint == "" {
		minioEndpoint = "minio:9000"
	}
	minioAccess = os.Getenv("MINIO_ACCESS_KEY")
	if minioAccess == "" {
		minioAccess = "minioadmin"
	}
	minioSecret = os.Getenv("MINIO_SECRET_KEY")
	if minioSecret == "" {
		minioSecret = "minioadmin"
	}

	var err error
	minioClient, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccess, minioSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("MinIO接続失敗: %s", err)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/images", imagesHandler)
	// 静的ファイル（HTML, CSS, JS）は ./static 配下に配置
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("APIサーバー起動 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
