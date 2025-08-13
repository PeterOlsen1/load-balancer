# The only 2 commands I'll need
# run 'docker ps' to see running docker containers
# run 'docker stop {id}' to stop a running container

FROM node:23

WORKDIR /app

COPY package.json package-lock.json ./

RUN npm install

RUN npm install -g typescript

COPY . .

EXPOSE 3000

CMD ["npm", "run", "build"]