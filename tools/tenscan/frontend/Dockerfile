FROM node:18-buster as runner
# Install pnpm
RUN npm install -g pnpm
# setup container data structure
RUN mkdir -p /home/obscuro/go-obscuro/tools/tenscan/
COPY ./tools/tenscan/frontend /home/obscuro/go-obscuro/tools/tenscan/frontend

WORKDIR /home/obscuro/go-obscuro/tools/tenscan/frontend
RUN pnpm install

EXPOSE 80