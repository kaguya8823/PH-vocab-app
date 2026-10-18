## PH-vocab-app
english-vocab-app をもとに、 prisma から golang に差し替えて作ります。

## Go をインストール
backend 内
go mod init english-vocab-app/backend
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql