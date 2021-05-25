FROM golang:alpine as build
ADD . /go/src/github.com/ms-choudhary/quotes
RUN go install github.com/ms-choudhary/quotes@latest

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/quotes /usr/bin/quotes
EXPOSE 9090
CMD ["quotes"]
