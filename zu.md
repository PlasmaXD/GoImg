# マークダウン上でPlantUMLを書く方法

VSCodeでマークダウンファイル内にPlantUMLを埋め込む方法をご紹介します。

## 必要な拡張機能

マークダウン内でPlantUMLを表示するには、以下の拡張機能をインストールします：

1. **Markdown Preview Enhanced** - マークダウン内のPlantUMLをレンダリングするための拡張機能
2. **PlantUML** - PlantUML構文のサポートと描画機能を提供する拡張機能

VSCodeの拡張機能マーケットプレイス（Ctrl+Shift+X）から検索してインストールできます。

## マークダウン内でPlantUMLを書く方法

### 方法1: コードブロックを使用する

マークダウンファイル内で以下のように記述します：

````markdown
```
@startuml
class User {
  +id: int
  +name: string
}
class Order {
  +id: int
  +user_id: int
}
User "1" -- "many" Order : places
@enduml
```
```

### 方法2: Markdown Preview Enhancedの構文を使用する

````markdown
```
@startuml
actor User
participant "First Class" as A
participant "Second Class" as B
User -> A: DoWork
activate A
A -> B: Create Request
activate B
B -> B: DoWork
return Result
deactivate B
A -> User: Display Result
deactivate A
@enduml
```
```

## プレビュー方法

1. マークダウンファイルを開いた状態で、右クリックして「Markdown Preview Enhanced: Open Preview」を選択
2. または、キーボードショートカット `Ctrl+K V` を使用

## 画像としてエクスポート

Markdown Preview Enhancedでプレビューを開いた状態で：

1. プレビュー画面を右クリック
2. 「PlantUML」→「PNG」（または他の形式）を選択

## 設定のカスタマイズ

VSCodeの設定（Ctrl+,）で以下の項目をカスタマイズできます：

```json
"markdown-preview-enhanced.plantumlJarPath": "/path/to/plantuml.jar",
"markdown-preview-enhanced.usePandocParser": true
```

## トラブルシューティング

プレビューが表示されない場合：

1. GraphvizとJavaが正しくインストールされているか確認
2. VSCodeを再起動
3. 拡張機能の設定でPlantUMLのパスが正しく設定されているか確認

## 高度な使用法

### 外部PlantUMLファイルのインクルード

```markdown
@import "path/to/diagram.puml"
```

### テーマの適用

```plantuml
@startuml
!theme cerulean
class User {
  +id: int
}
@enduml
```

これらの方法を使えば、マークダウン内でPlantUMLダイアグラムを効率的に作成・表示できます。

---
Perplexity の Eliot より: pplx.ai/share