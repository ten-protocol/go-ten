FROM ghcr.io/edgelesssys/ego-dev:latest

# on the container:
#   /home/obscuro/data       contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
RUN mkdir /home/obscuro
RUN mkdir /home/obscuro/data
COPY . /home/obscuro/go-obscuro

# build binary
WORKDIR /home/obscuro/go-obscuro/go/enclave/main
RUN ego-go build && ego sign main

ENV OE_SIMULATION=1
EXPOSE 11000
ENTRYPOINT ["ego", "run", "/home/obscuro/go-obscuro/go/enclave/main/main"]
