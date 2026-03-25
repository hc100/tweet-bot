# tweet-bot

X に自動投稿するためのシンプルな Go 製 daemon です。

## 概要

- 毎朝 6:00 に `もにゅん` を投稿します
- 毎日 10:00 から 22:00 の間で 1 日 1 回、`2027年まで○日○時間○分○秒です。` を投稿します
- 投稿仕様は job を追加するだけで拡張できます

カウントダウン投稿の時刻は、日付ごとに一意に決まる擬似ランダム値です。同じ日なら再起動しても同じ時刻になります。

## 必要な環境変数

`.env` もしくは実環境の環境変数から、以下を読み込みます。

- `X_API_KEY`
- `X_API_SECRET`
- `X_ACCESS_TOKEN`
- `X_ACCESS_TOKEN_SECRET`
- `TZ` 任意。未指定時は `Asia/Tokyo`

`.env` の例:

```dotenv
X_API_KEY=...
X_API_SECRET=...
X_ACCESS_TOKEN=...
X_ACCESS_TOKEN_SECRET=...
TZ=Asia/Tokyo
```

## 起動方法

```bash
cp .env.example .env
go run ./cmd/tweet-bot
```

## 構成

- `cmd/tweet-bot`: エントリポイント
- `internal/config`: `.env` と環境変数の読み込み
- `internal/xclient`: X API への投稿
- `internal/scheduler`: 常駐スケジューラ
- `internal/jobs`: 投稿仕様ごとの job 定義

## job を追加する方法

`internal/scheduler.Job` を満たす型を追加し、`cmd/tweet-bot/main.go` に登録してください。

必要なメソッド:

```go
type Job interface {
	Name() string
	NextRun(after time.Time) (time.Time, bool)
	BuildPost(now time.Time) (string, error)
}
```

## systemd での運用例

Ubuntu での unit file 例です。

```ini
[Unit]
Description=tweet-bot
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
WorkingDirectory=/opt/tweet-bot
ExecStart=/opt/tweet-bot/tweet-bot
Restart=always
RestartSec=5
EnvironmentFile=/opt/tweet-bot/.env

[Install]
WantedBy=multi-user.target
```

ビルド例:

```bash
GOOS=linux GOARCH=amd64 go build -o tweet-bot ./cmd/tweet-bot
```

## テスト

```bash
go test ./...
```
