FROM golang:1.20

ENV TODO_PORT :8000
ENV TODO_DBFILE db/scheduler.db
ENV TODO_PASSWORD TODO_PASSWORD
ENV SECRET_KEY @F#SJ;9QgMOAGxY33.0d93r(_:;nS:

WORKDIR /my_app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8000
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /my_app

CMD ["/my_app"]