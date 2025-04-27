FROM golang:1.24 AS build
COPY ./ /go/src/unchained
WORKDIR /go/src/unchained
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -v -o /unchained/unchained -tags with_reality_server ./cmd/unchained

FROM alpine
COPY --from=build /unchained/unchained /unchained/unchained

WORKDIR /unchained/
ENV UNCHAINED_DIR=/data/

ENTRYPOINT [ "/unchained/unchained" ]
CMD [ "run" ]