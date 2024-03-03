FROM golang:alpine
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build updateRedis.go 
RUN go build simulatorApp.go
RUN chmod +x run.sh
CMD ["./run.sh"]