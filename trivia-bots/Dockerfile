FROM node:14.15.4

# create destination directory
RUN mkdir -p /trivia-bots
WORKDIR /trivia-bots

# copy the app, note .dockerignore
COPY . /trivia-bots
RUN npm install

RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get install -y ffmpeg

# start the app
CMD [ "node", "index.js" ]