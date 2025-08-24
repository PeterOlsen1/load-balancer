# The only 2 commands I'll need
# run 'docker ps' to see running docker containers
# run 'docker stop {id}' to stop a running container

FROM node:23-alpine

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

RUN npm install -g typescript

COPY . .

RUN npm run build

EXPOSE 3000

CMD ["node", "build/server.js"]