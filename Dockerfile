FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o exparams

EXPOSE 8091

CMD ["./exparams"]