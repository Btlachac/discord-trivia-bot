FROM node:latest

ENV NODE_OPTIONS=--openssl-legacy-provider

# create destination directory
RUN mkdir -p /trivia-ui
WORKDIR /trivia-ui

# ARG API_URL
# ENV API_URL=$API_URL

# copy the app, note .dockerignore
COPY . /trivia-ui
RUN npm install

# build necessary, even if no static files are needed,
# since it builds the server as well
RUN npm run build

# expose 5000 on container
EXPOSE 5000

# set app serving to permissive / assigned
ENV NUXT_HOST=0.0.0.0
# set app port
ENV NUXT_PORT=5000

# start the app
CMD [ "npm", "start" ] 