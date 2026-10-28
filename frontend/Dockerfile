## フロントエンド用の Dockerfile の構造
FROM node:24.22.1

WORKDIR /src

RUN npm install --global npm@11.19.1

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]