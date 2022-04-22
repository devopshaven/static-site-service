FROM golang:1.18 as build

LABEL org.opencontainers.image.source = "https://github.com/devopshaven/static-site-service"
LABEL maintainer="Gyula Paal <paalgyula@paalgyula.com>"

ENV GOPROXY=goproxy.pirat.app

WORKDIR /usr/src/app
COPY go.mod go.sum ./

CMD go mod download -x

COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -ldflags="-s -w" -o /app .

# Alpine latest
FROM alpine:latest
COPY --from=build /app /app

CMD ["/app"]
