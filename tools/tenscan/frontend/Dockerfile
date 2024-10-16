FROM node:18-buster as runner

RUN npm install -g pnpm

COPY . /home/obscuro/go-obscuro/

WORKDIR /home/obscuro/go-obscuro/

RUN pnpm install

WORKDIR /home/obscuro/go-obscuro/tools/tenscan/frontend

RUN pnpm install --filter ./tools/tenscan/frontend...

EXPOSE 80

CMD ["pnpm", "run", "start"]
