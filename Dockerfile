FROM golang:1.24 AS build
COPY ./ /go/src/unchained
WORKDIR /go/src/unchained
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -v -o /bins/unchained ./cmd/unchained

FROM alpine
COPY --from=build /bins/unchained /bins/unchained

ENTRYPOINT [ "/bins/unchained"]
CMD [  "run" ]