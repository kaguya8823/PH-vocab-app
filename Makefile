.PHONY: help doctor deps up down build restart logs ps \
        backend-shell db db-shell frontend-dev \
        test-backend lint-frontend clean

# デフォルトは help を表示
.DEFAULT_GOAL := help

help: ## コマンド一覧を表示
	@awk 'match($$0, /^[a-zA-Z_-]+:.*## /) { \
		target = substr($$0, 1, index($$0, ":") - 1); \
		description = substr($$0, RSTART + RLENGTH); \
		printf "  \033[36m%-16s\033[0m %s\n", target, description \
	}' $(MAKEFILE_LIST)

doctor: ## 開発に必要なツールと Docker 接続を確認
	@missing=0; \
	for cmd in make docker go node npm; do \
		if command -v "$$cmd" >/dev/null 2>&1; then \
			printf 'OK: %s (%s)\n' "$$cmd" "$$(command -v "$$cmd")"; \
		else \
			printf '不足: %s がインストールされていません\n' "$$cmd"; \
			missing=1; \
		fi; \
	done; \
	if command -v docker >/dev/null 2>&1; then \
		if docker compose version >/dev/null 2>&1; then \
			printf 'OK: Docker Compose v2\n'; \
		else \
			printf '不足: Docker Compose v2 プラグインがありません\n'; \
			missing=1; \
		fi; \
		if docker info >/dev/null 2>&1; then \
			printf 'OK: Docker デーモンに接続できます\n'; \
		else \
			printf '不足: Docker デーモンに接続できません (起動状態と権限を確認してください)\n'; \
			missing=1; \
		fi; \
	fi; \
	if [ "$$missing" -ne 0 ]; then \
		printf 'README.md の「開発環境の準備」を確認してください\n'; \
		exit 1; \
	fi

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

db-shell: db ## MySQLに接続 (db の別名)

# ---------- 依存関係 ----------
deps: ## Go とフロントエンドの依存関係を取得
	cd backend && go mod download
	cd frontend && npm ci

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