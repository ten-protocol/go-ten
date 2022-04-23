FROM ghcr.io/edgelesssys/ego-dev:latest

RUN git clone https://github.com/obscuronet/obscuro-playground
RUN cd obscuro-playground/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "obscuro-playground/go/obscuronode/enclave/main/main"]
EXPOSE 11000
