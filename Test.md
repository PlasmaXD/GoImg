
# Mermaid 



### 1. アーキテクチャ図（Mermaid 記法）

```mermaid
graph LR
    A[ブラウザ (ユーザーUI)] -->|画像アップロード| B[APIサーバー]
    B -->|ジョブ送信| C[RabbitMQ]
    C -->|ジョブ配信| D[Worker]
    D -->|画像処理・アップロード| E[MinIO]
```

処理の流れ

---

### 2. クラス仕様図（Mermaid 記法）

```mermaid
classDiagram
    class APIServer {
      +uploadHandler(w, r)
      +main()
      -amqpURL : string
      -queueName : string
    }
    class Worker {
      +main()
      +connectRabbitMQ()
      +processImage(msg)
      +uploadToMinIO(imageData)
      -amqpURL : string
      -queueName : string
      -minioEndpoint : string
      -minioAccess : string
      -minioSecret : string
      -minioBucket : string
    }
```

APIサーバーと Worker の主なメソッドとフィールド

---

### 3. 処理フロー図（Mermaid 記法）

```mermaid
sequenceDiagram
    participant U as ユーザー
    participant A as APIサーバー
    participant R as RabbitMQ
    participant W as Worker
    participant M as MinIO

    U->>A: 画像アップロード
    A->>U: 「画像受け付け」メッセージ
    A->>R: 画像データを送信
    R->>W: ジョブ配信
    W->>W: 画像処理（リサイズ・サムネイル生成）
    W->>M: 処理済み画像をアップロード
```

一連の流れ(シーケンス図)
---

Merkdown 内に記述した Mermaid コードをレンダリングできます。  
- **オンラインエディタ**: [Mermaid Live Editor](https://mermaid.live/) 