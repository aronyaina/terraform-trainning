FROM golang:1.21

WORKDIR /app
COPY go.mod go.sum ./

COPY vendor ./
COPY index.html ./
RUN go mod download
COPY *.go ./

EXPOSE 8000
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-fibonnaci

CMD ["/docker-fibonnaci"]