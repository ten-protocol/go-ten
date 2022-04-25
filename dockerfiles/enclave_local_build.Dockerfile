FROM ghcr.io/edgelesssys/ego-dev:latest

# build the enclave from the current branch
RUN mkdir /home/obscuro-playground
COPY ./ /home/obscuro-playground
RUN ls /home/obscuro-playground  && cd /home/obscuro-playground/go/obscuronode/enclave/main && ego-go build && ego sign main

ENV OE_SIMULATION=1
ENTRYPOINT ["ego", "run", "obscuro-playground/go/obscuronode/enclave/main/main"]
EXPOSE 11000
