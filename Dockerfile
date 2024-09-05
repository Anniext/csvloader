FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV GOCACHE=/root/.cache/go-build
ENV GOPATH=/root/go
ENV GOBIN=$GOPATH/bin
ENV PATH=$PATH:$GOBIN

RUN apk add --no-cache git

RUN make build

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/build/csvloader /usr/local/bin/csvloader

CMD ["csvloader"]

# docker build -t my-go-app .
# docker run -it --rm --name xxxx csvloader