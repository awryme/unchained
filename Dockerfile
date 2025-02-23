FROM golang:1.24 AS build
COPY ./ /go/src/unchained
WORKDIR /go/src/unchained
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -v -o /bins/unchained -tags with_reality_server ./cmd/unchained

FROM alpine
COPY --from=build /bins/unchained /bins/unchained

WORKDIR /config

ENTRYPOINT [ "/bins/unchained", "-c", "/config/unchained.json" ]
CMD [  "run" ]