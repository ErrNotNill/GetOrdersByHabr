FROM golang:latest

ENV GOPATH=/
RUN go env -w GO111MODULE=auto

COPY ./ ./

RUN go mod download

#make database dump

RUN go build -o search_panel ./cmd/main.go

CMD ["./search_panel"]