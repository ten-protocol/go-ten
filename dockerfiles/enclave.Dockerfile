FROM ghcr.io/edgelesssys/ego-dev:latest

# build the enclave from the current branch
RUN mkdir /home/obscuro-playground
COPY . /home/obscuro-playground
RUN cd /home/obscuro-playground/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
EXPOSE 11000
ENTRYPOINT ["ego", "run", "/home/obscuro-playground/go/obscuronode/enclave/main/main"]
