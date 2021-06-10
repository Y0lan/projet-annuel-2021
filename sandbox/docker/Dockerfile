FROM bitnami/minideb
MAINTAINER YOLAN MALDONADO / LUKA BOELS

RUN apt-get update

# install python3
RUN install_packages software-properties-common

# install go
RUN install_packages golang-go

# install rust
RUN install_packages curl build-essential
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y

# install isolate
RUN install_packages libcap-dev git
RUN git clone https://github.com/ioi/isolate /tmp/isolate
RUN cd /tmp/isolate && make isolate && make install
RUN mv /tmp/isolate/isolate /bin/isolate

# install tests
COPY tests /tests

