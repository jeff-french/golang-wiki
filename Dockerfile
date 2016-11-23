FROM golang:1.7

COPY . /go/src/app

WORKDIR /go/src/app

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]

