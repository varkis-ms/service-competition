FROM golang:1.22
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go mod tidy
RUN go build ./cmd/service
ENTRYPOINT /app/service
