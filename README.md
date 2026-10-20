## PH-vocab-app
english-vocab-app をもとに、 prisma から golang に差し替えて作ります。

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

mysql の見方
https://claude.ai/share/4288967e-b54e-4b52-8fa3-6c82d62d7cdf

## ポートについて
必要があれば、ポートの追加で5174を手動でしなければならない。

## docker 起動と停止
起動：docker compose up -d

停止：docker compose down -v