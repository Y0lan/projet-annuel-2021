FROM bitnami/minideb
MAINTAINER YOLAN MALDONADO / LUKA BOELS

ENV HOME /home/runner
ENV API $HOME/api
ENV USER runner



# prepare the system
RUN apt-get update
RUN apt-get upgrade -y

# install python3
RUN install_packages software-properties-common

# install go
RUN install_packages golang-go git

# install rust dependencies
RUN install_packages curl build-essential

# create user
RUN useradd -c "$USER user" -m -d $HOME -s /bin/bash $USER
RUN mkdir $API
RUN chown -R $USER $API

# install rust as user
USER $USER
WORKDIR $HOME
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y
ENV PATH $PATH:$HOME/.cargo/bin

# add project files
COPY compile.go $API/compile.go
ADD tests $API/tests


# verify everything is working
WORKDIR $API
RUN ./tests/test_all.sh

# get go dependencies for the api
RUN go get -d ./...

# build the api
RUN go build -o main .

# run the api
CMD ["/home/runner/api/main"]

