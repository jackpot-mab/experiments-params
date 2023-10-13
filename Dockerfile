FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o experiments-params

EXPOSE 8091

CMD ["./experiments-params"]