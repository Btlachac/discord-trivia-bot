FROM node:14.15.4

RUN mkdir /trivia-ui
WORKDIR /trivia-ui

RUN npm install

EXPOSE 5000

ENV NUXT_HOST=0.0.0.0
ENV NUXT_PORT=5000

CMD [ "npm", "run", "dev" ]