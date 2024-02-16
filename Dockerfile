FROM golang:latest

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build .

EXPOSE 8080

ENTRYPOINT ["./task-manager-golang", "-environment=fly"]
