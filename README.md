## PH-vocab-app
english-vocab-app をもとに、 prisma から golang に差し替えて作ります。
単語や例文に録音機能を追加して、録画・再生・更新もできるようにします。

## Go をインストール
backend 内
go mod init english-vocab-app/backend
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql

## frontend
ルートディレクトリで
npm create vite@latest frontend
 -> React , TypeScript , ESLint 

frontend 内
npm install
npm run dev で起動

npm install axios

## mysql
docker exec -it PH_mysql mysql -uroot -p PH_vocab

## mysql に、音声を録音、保存、再生するテーブルを作る
上記のコマンドでmysql に入る
show tables; で中身を見る。 作った状態では vocabularies のみ

mysql> CREATE DATABASE vocabulary_audios; と入力

mysql> show databases; で中をのぞくと、上記で作ったデータベースが作成されている。
このデータベースの中に、音声を録音や保存、再生できる処理を作る。


mysql の見方
https://claude.ai/share/4288967e-b54e-4b52-8fa3-6c82d62d7cdf

## ポートについて
必要があれば、ポートの追加で5174を手動でしなければならない。

## docker 起動と停止
起動：docker compose up -d

停止：docker compose down -v

再起動：docker compose up -d --build backend frontend

再ビルド：docker compose up -d --build