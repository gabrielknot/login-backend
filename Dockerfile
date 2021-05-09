WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM scratch
COPY --from=builder /app/httpserver /app/
EXPOSE 3080
ENTRYPOINT ["/app/httpserver"]
FROM golang











