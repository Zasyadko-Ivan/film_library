FROM golang:latest

RUN apt update && apt install -y zsh && apt install -y nano

ENTRYPOINT ["/bin/zsh"]
