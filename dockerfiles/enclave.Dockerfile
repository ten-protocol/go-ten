FROM ghcr.io/edgelesssys/ego-dev:latest

RUN git clone https://github.com/obscuronet/go-obscuro
RUN cd go-obscuro/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "go-obscuro/go/obscuronode/enclave/main/main"]
EXPOSE 11000
