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