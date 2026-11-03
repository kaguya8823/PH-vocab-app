# PH-vocab-app

Go / MySQL の API と React / Vite のフロントエンドで構成された単語学習アプリです。単語・例文の音声録音、保存、再生にも対応しています。

## 開発環境の準備

必要なツール:

- GNU Make
- Docker Engine と Docker Compose v2 (Docker デーモンを起動し、現在のユーザーから接続できること)
- Go (バージョンは `backend/go.mod` を参照)
- Node.js と npm (Vite が対応する Node.js バージョンが必要)

Ubuntu では、Make / Go / Node.js / npm / Compose プラグインを次のように導入できます。

```sh
sudo apt update
sudo apt install make golang-go nodejs npm docker-compose-v2
```

Docker Engine 自体がない場合は、利用中の Ubuntu 向け Docker Engine のインストール手順に従ってください。Docker を起動した後も権限エラーになる場合は、ユーザーを `docker` グループに追加して再ログインするか、管理者に設定を依頼してください。

導入状態はリポジトリのルートで確認できます。

```sh
make doctor
```

## 起動

Docker Compose で MySQL・API・フロントエンドをまとめて起動します。初回はイメージのビルドに時間がかかります。API が MySQL にテーブルを作成するため、DB の手動セットアップは不要です。

```sh
make deps      # Go モジュールと npm パッケージを取得 (ローカル開発用)
make up        # コンテナをビルドしてバックグラウンド起動
make logs      # ログを確認 (終了は Ctrl+C)
make ps        # コンテナ状態を確認
```

`make up` はコンテナがない場合にイメージもビルドします。ソース変更後など明示的に再ビルドする場合は `make build` の後に `make up` を実行してください。

- アプリ: http://localhost:5174
- API: http://localhost:3003
- MySQL: `localhost:3308` (Compose 内のポートは 3306)

## ローカルで個別に開発

依存関係を取得した後、MySQL は Compose で起動し、フロントエンドやバックエンドをホスト上で個別に実行できます。

```sh
docker compose up -d mysql
```

バックエンドには `DB_DSN` が必要です。Compose と同じ既定の DB 設定で起動する場合:

```sh
cd backend
DB_DSN='root:root@tcp(127.0.0.1:3308)/PH_vocab?parseTime=true' go run ./cmd/server
```

別のターミナルでフロントエンドを起動します。

```sh
make frontend-dev
```

Vite は `/api` を `http://localhost:3003` に転送します。別ポートの API を使う場合は `VITE_API_PROXY_TARGET` を設定してください。

## よく使うコマンド

```sh
make help
make test-backend
make lint-frontend
make down
make db        # MySQL に接続
```

`make down` はコンテナを停止・削除しますが、DB データは保持します。`make clean` は Docker ボリュームも削除するため、MySQL のデータと音声ファイルが消えます。
