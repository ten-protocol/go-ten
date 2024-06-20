FROM node:18-buster as runner
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/contracts/src/bridge
COPY ./contracts/src/bridge/frontend /home/obscuro/go-obscuro/contracts/src/bridge/frontend

WORKDIR /home/obscuro/go-obscuro/contracts/src/bridge/frontend
RUN npm install

EXPOSE 80