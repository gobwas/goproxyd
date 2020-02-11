FROM golang:1.13.6-alpine

WORKDIR /go/src/goproxy

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go install -v .

EXPOSE 8080
CMD ["goproxyd"]
