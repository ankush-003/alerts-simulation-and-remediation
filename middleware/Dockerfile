FROM golang:alpine
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go run createAlerts.go 
RUN go build demoSimulator.go
CMD ["./demoSimulator"]