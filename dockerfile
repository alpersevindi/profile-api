FROM golang:1.20

WORKDIR /go/src/app

COPY . .

RUN go get -u github.com/labstack/echo/v4
RUN go get -u github.com/aws/aws-sdk-go
RUN go get -u github.com/google/uuid

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
