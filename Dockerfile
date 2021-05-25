FROM golang
ADD . /go/src/github.com/ms-choudhary/quotes
RUN go install github.com/ms-choudhary/quotes@latest
EXPOSE 9090
ENTRYPOINT ["quotes"]
