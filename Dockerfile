FROM golang:1.23-alpine

WORKDIR /dockerapp

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY .env ./
COPY pb_data /dockerapp/pb_data
COPY handler /dockerapp/handler
COPY helper /dockerapp/helper
COPY middleware /dockerapp/middleware
COPY model /dockerapp/model

RUN go build -o /go-exec

EXPOSE 8080

CMD ["/go-exec", "serve", "--http=0.0.0.0:8080"]
