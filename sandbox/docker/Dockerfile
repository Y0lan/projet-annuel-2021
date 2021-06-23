FROM bitnami/minideb
MAINTAINER YOLAN MALDONADO / LUKA BOELS

# update
RUN apt-get update
RUN apt-get upgrade -y

# install python3
RUN install_packages software-properties-common

# install go
RUN install_packages golang-go

# install rust dependencies
RUN install_packages curl build-essential

# create user to not execute code as root
RUN useradd -c "unprivileged user" -m -d /home/unprivileged -s /bin/bash unprivileged
USER unprivileged
ENV HOME /home/unprivileged

# install rust
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH $PATH:$HOME/.cargo/bin


HEALTHCHECK CMD ./$HOME/shared/tests/test_all.sh || exit 1
COPY compile.go compile.go
RUN go run compile.go

# install isolate
# RUN install_packages libcap-dev git
# RUN git clone https://github.com/ioi/isolate /tmp/isolate
# RUN cd /tmp/isolate && make isolate && make install
# RUN mv /tmp/isolate/isolate /bin/isolate


