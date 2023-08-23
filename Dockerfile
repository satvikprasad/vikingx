FROM node:lts-alpine

WORKDIR /app/frontend
COPY ./frontend/package*.json ./

RUN npm install

COPY ./frontend ./

RUN npm run build

FROM golang:latest

WORKDIR /app
COPY . /app

RUN go mod download

RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]

