FROM golang:1.17.1 as golang

WORKDIR /howareyou
COPY . .

WORKDIR /howareyou/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/howareyou


FROM gcr.io/distroless/base

COPY --from=golang /go/bin /app
CMD ["app/howareyou"]
