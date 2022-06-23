FROM ghcr.io/edgelesssys/ego-dev:latest

# on the container:
#   /home/obscuro/data       contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
RUN mkdir /home/obscuro
RUN mkdir /home/obscuro/data
WORKDIR /home/obscuro

RUN git clone https://github.com/obscuronet/go-obscuro
RUN cd go-obscuro/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "/home/obscuro/go-obscuro/go/obscuronode/enclave/main/main"]