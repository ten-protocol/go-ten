FROM node:18-buster as runner
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/tools/obscuroscan_v2/
COPY ./tools/obscuroscan_v2/frontend /home/obscuro/go-obscuro/tools/obscuroscan_v2/frontend

WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan_v2/frontend
RUN npm install && npm install http-server -g

WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan_v2/frontend
EXPOSE 80