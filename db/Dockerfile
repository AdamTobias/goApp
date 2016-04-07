FROM golang:1.6

ADD . /go/src/github.com/adamtobias/goApp/db
RUN go get github.com/go-sql-driver/mysql 
RUN go install github.com/adamtobias/goApp/db

ENTRYPOINT env=docker db
# CMD ["dbLayer", "env=docker"]

EXPOSE 8000