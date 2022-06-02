FROM ghcr.io/edgelesssys/ego-dev:latest

# build the enclave from the current branch
RUN mkdir /home/go-obscuro
COPY . /home/go-obscuro
RUN cd /home/go-obscuro/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "/home/go-obscuro/go/obscuronode/enclave/main/main"]
EXPOSE 11000
