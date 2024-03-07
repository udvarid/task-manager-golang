FROM ubuntu:latest

RUN apt update && apt install -y ca-certificates

COPY task-manager-golang ./
COPY ./templates/*.html ./templates/
RUN mkdir /db

EXPOSE 8080

ENTRYPOINT ["./task-manager-golang", "-environment=fly"]
