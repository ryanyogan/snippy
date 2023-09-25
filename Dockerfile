FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
COPY ui/ ui/

RUN go mod download

COPY . .

RUN go build -o app-server ./cmd/web

EXPOSE 8080

ENTRYPOINT [ "./app-server" ]
