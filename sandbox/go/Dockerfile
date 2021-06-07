# Sandbox part
FROM debian:buster as sandbox

ENV GOPATH /go
ENV GOCACHE /go/.cache
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV GOROOT_BOOTSTRAP /usr/local/gobootstrap
ENV GO_VERSION 1.16

RUN apt-get update && apt-get install -y ca-certificates curl

# Install GO
RUN curl -sSL https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz -o /tmp/go.tar.gz
RUN curl -sSL https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz.sha256 -o /tmp/go.tar.gz.sha256
RUN echo "$(cat /tmp/go.tar.gz.sha256) /tmp/go.tar.gz" | sha256sum -c -
RUN tar -C /usr/local/ -vxzf /tmp/go.tar.gz

# Add and compile sandbox
COPY . /go/src/sandbox
RUN go install sandbox

# Running user input

FROM debian:buster as executor

RUN apt-get update && apt-get install -y ca-certificates git --no-install-recommends

COPY --from=builder /usr/local/go /usr/local/go

ENV GOPATH /go
ENV GOCACHE /tmp/.cache

ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH

RUN mkdir /app
COPY --from=sandbox /go/bin/sandbox /app
WORKDIR /app

# Run tests
RUN /app/sandbox test

EXPOSE 8080
ENTRYPOINT ["/app/sandbox"]

