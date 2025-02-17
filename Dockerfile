FROM golang:1.22.5

ENV TODO_PORT :8000
ENV TODO_DBFILE storage/scheduler.db
ENV TODO_PASSWORD TODO_PASSWORD

WORKDIR /my_app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app

CMD ["/my_app"]