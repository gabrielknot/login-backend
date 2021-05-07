FROM golang:1.10

WORKDIR /
COPY . .
RUN go get -d github.com/gorilla/mux
RUN go get -d github.com/rs/cors
RUN go get -d github.com/go-sql-driver/mysql

CMD ["go","run","main.go"]
