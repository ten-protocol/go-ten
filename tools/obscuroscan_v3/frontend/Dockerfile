FROM node:18-buster as runner
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/tools/obscuroscan_v3/
COPY ./tools/obscuroscan_v3/frontend /home/obscuro/go-obscuro/tools/obscuroscan_v3/frontend

WORKDIR /home/obscuro/go-obscuro/tools/obscuroscan_v3/frontend
RUN npm install && npm install http-server -g

EXPOSE 80