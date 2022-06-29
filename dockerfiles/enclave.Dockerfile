FROM ghcr.io/edgelesssys/ego-dev:latest

# on the container:
#   /home/obscuro/data       contains working files for the enclave
#   /home/obscuro/go-obscuro contains the src
RUN mkdir /home/obscuro
RUN mkdir /home/obscuro/data
COPY . /home/obscuro/go-obscuro
RUN cd /home/obscuro/go-obscuro/go/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "/home/go-obscuro/go/enclave/main/main"]
