# Trivia Server

This part of the project is an API used primarily for inserting new trivia from the UI into the Postgres database, retrieving unused trivia and marking trivia as used for the bots.

## Setting up Postgres

This project is set up to work with a Postgres database. <br>

We utilize [golang-migrate](https://github.com/golang-migrate/migrate) for automatic database migrations on server start.<br>

All you'll need to do is set up an empty postgres database, point the server at it, and the migrations to set up the DB will run when you start the server.

## Installation

### Prerequisites

You will need to have Go installed to run this project. <br>

### Environment Variables

You can find an example.env file in this folder. Rename the file to '.env' and replace the placeholders with your environment variables.

#### How to Run

Build project:<br>
```
go build
```

Start the application:<br>
```
./go-trivia-api
```
