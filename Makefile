.PHONY: help up down build restart logs ps \
        backend-shell db-shell frontend-dev \
        test-backend lint-frontend clean

# デフォルトは help を表示
.DEFAULT_GOAL := help

help: ## コマンド一覧を表示
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ---------- Docker ----------
up: ## コンテナをバックグラウンドで起動
	docker compose up -d

down: ## コンテナを停止・削除
	docker compose down

build: ## イメージを再ビルド
	docker compose build --no-cache

restart: ## 再起動
	docker compose down
	docker compose up -d

logs: ## ログを追跡表示
	docker compose logs -f

ps: ## 起動中のコンテナ一覧
	docker compose ps

# ---------- 各コンテナに入る ----------
backend-shell: ## Goコンテナのシェルに入る
	docker compose exec backend sh

db: ## MySQLに接続
	docker exec -it PH_mysql mysql -uroot -p PH_vocab

# ---------- フロントエンド ----------
frontend-dev: ## React開発サーバーを起動
	cd frontend && npm run dev

lint-frontend: ## ESLintを実行
	cd frontend && npm run lint

# ---------- バックエンド ----------
test-backend: ## Goのテストを実行
	cd backend && go test ./...

# ---------- 掃除 ----------
clean: ## コンテナ・ボリュームを完全削除(DBデータも消えます)
	docker compose down -v