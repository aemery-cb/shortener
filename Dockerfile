FROM node:18-alpine as ui 

WORKDIR /app

COPY package*.json ./
COPY . .
RUN npm i 

RUN npm run build

FROM golang:1.18 as server

WORKDIR /app
COPY go.* ./
COPY *.go ./
COPY --from=ui /app/public ./public
RUN go build

FROM alpine

COPY --from=server /app/shortener /

CMD [/shortener]


