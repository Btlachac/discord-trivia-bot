FROM node:14.15.4

# create destination directory
RUN mkdir -p /trivia-bots
WORKDIR /trivia-bots

# copy the app, note .dockerignore
COPY . /trivia-bots
RUN npm install

# start the app
CMD [ "node", "index.js" ]