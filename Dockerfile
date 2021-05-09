FROM golang as builder
WORKDIR /app
ENV FRONT_PORT=80
ENV HOST="localhost"
ENV MYSQL_PORT="52001"
ENV MYSQL_USER="root"
ENV MYSQL_PASSWORD="admin"
ENV MYSQL_DATABASE="usersDb"
COPY . .
RUN go mod download

COPY go.mod .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .
CMD ["go","run","main.go"]
#FROM scratch
#COPY --from=builder /app/httpserver /app/
#EXPOSE 3080
#ENTRYPOINT ["/app/httpserver"]
#
#
#
#
#
#
#
#
#
#
#
