FROM node:18-buster as runner
# setup container data structure
RUN mkdir -p /home/obscuro/go-ten/tools/contracts/src/bridge
COPY ./contracts/src/bridge/frontend /home/obscuro/go-ten/contracts/src/bridge/frontend

WORKDIR /home/obscuro/go-ten/tools/tenscan/frontend
RUN npm install

EXPOSE 80