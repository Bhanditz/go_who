FROM golang:1.11.4-alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache bash git

RUN go get -u github.com/gorilla/mux && \
    go get -u github.com/gorilla/rpc && \
    go get -u github.com/gorilla/sessions && \    
    go get -u github.com/micro/protobuf/protoc-gen-go && \
    go get -u github.com/mchirico/gog/cmd/gog  && \
    go get -u github.com/dgrijalva/jwt-go && \
    go get -u github.com/levigross/grequests && \    
    go get -u github.com/mchirico/go_who/... && \
    go get -u github.com/prometheus/node_exporter

RUN go build -o /auth github.com/mchirico/go_who/cmd/auth
RUN cd /go/src/github.com/prometheus/node_exporter && go build



FROM golang:1.11.4-alpine

COPY --from=builder /go/src/github.com/prometheus/node_exporter/node_exporter /go/bin/node_exporter
COPY --from=builder /auth /go/bin/auth

COPY exec.go /go/src/exec.go
COPY .go_who_secret_google /root/.go_who_secret_google

CMD ["/usr/local/go/bin/go","run","/go/src/exec.go"]



# To build this:
# docker build -t us.gcr.io/mchirico/auth:pig .


# To run this:
# docker run -p 4591:4591 -p 9100:9100 --rm -it us.gcr.io/mchirico/auth:pig


# To push
# docker push us.gcr.io/mchirico/auth:pig
