// # errors 
// # maybe switch to sqlx?
// # contexts and tracing ?
// # logs ?
// # onshutdown
// # swagger
// # migrations
// # switch from dev to prod

// dockerize
// users

// I did not add the ability to run multiple servers since it would be too much for a template that tries to be clean and clear
// but if you need you can surely add it without modifying much
// or other specific features
// this is supposed to give you a solid base from which you can add your specific features

// logs
// The logs directory path has been set to this root directory, but do not keep them here. Add your own path to a
// directory outside of the root project, like /var/log/myapp for Unix/Linux systems, where you can separate the concerns
// or create log rotation if needed.

// Set up for production is to add the variables in `~/.bash_profile` or `/etc/environment` or directly in terminal.

// I added two examples of outputting logs to file in /books and /book/{id}


// go install github.com/swaggo/swag/cmd/swag@latest
// go install github.com/cosmtrek/air@latest


# Golang Rest API Template
> Golang Rest Api Template with clear, scalable structure that can sustain large APIs.

## Table of Contents

- [Features](#Features)
- [Directory Structure](#Directory-Structure)
- [Description](#Description)
- [Setup](#Setup)
- [Template Tour](#Tempalte-Tour)
- [License](#license)

## Features

- swagger
- sqlx
- standard responses for success and fail
- errors standard
- loggers
- migrations
- hot reload

---

## Directory Structure
```
- /cmd --> Contains the app's entry-points 
  |- /server
     |- /docs
     |- main.go
     |- Makefile
  |- /another_binary
- /config --> Contains the config structures that the server uses.
- internal --> Contains the app's code
   |- /errors
   |- /handlers
   |- /middleware
   |- /model
   |- /storage
   |- server.go
- /logs --> The folder's files are not in version control. Here you'll have the logs of the apps (the ones you specify to be external)
- /migrations --> Migrations to set up the database schema on your db.
- /pkg --> Packages used in /internal
   |- /httputils
   |- /logger
- .air.toml
- .env --> Not in version control. Need to create your own - see below.
- .gitignore
- docker-compose.yml
- Dockerfile
- go.mod
- go.sum
- LICENSE
- README.md
```

## Description

**The Why**

I have spent a while looking for Go Api templates to automate my workflow but I had some trouble with the ones I've found. Some of them were way to specific, sometimes allowing just one set of handlers or one model in db to exist by default, in which case I had to rewrite and restructure it to make it open for extension, to add more handlers, DBs etc. Others were way to complex, with a deep hierarchical structure that didn't made sense (to me at least), half of which I would delete afterwards. So I wanted something that is as flat as possible but still having some structure that is easily extendable when needed, and also has the basic functionality set up so that I only need to add to it. This template is my attempt at achieving that.

To keep it simple, the template creates CRUD of books which then can be deleted for your specific handlers, but they show the way the api works generally. The server uses `gorilla/mux` as a router, `urfave/negroni` as a base middleware, `sirupsen/logrus` as a logger and `unrolled/render` as the functionality to format the responses in any way you want before automatically sending them. The API also uses `unrolled/secure` to improve the security. For database management `jmoiron/sqlx` is used to improve the ease of use.

The app uses Air as a hot-reloader and Docker for when you need to deploy it. You can start it in both ways (see [Setup](#setup)).
For environment variables `joho/godotenv` has been used instead of a `.yaml` file for security considerations. Again, everything is extendable, so you can add a `.yaml` file if you need more hierarchical structure or your environment variables don't need to be secure.
If you decide to not use Docker, in dev mode you use the variables from .env, and in production you add them in the terminal or in `~/.bash_profile`/`/etc/environment`. If you do decide to use Docker, keep the variables in .env in development mode, and add another .env on your production server with the prod variables in .env. The .env file is not version controlled so the workflow will be smooth. I choose to go with this approach from hundreds because I tried to hit the middle - simple and relatively secure, in a way that *most* people will use this template. If the api becomes large and you have specific needs for your case, you can add the variables in prod in command-line, in volumes, in a secret manager, in your kubernetes/docker swarm or any other way you want/need. But for most people and most cases, this approach will be more than enough.

I made sure you can add to it and modify without any pain, so for example, you can add one more db without modifying anything from the existing code, and also you can change the current db (Postgres) with any other also by not modifying anything besides creating the connect function for your new db. Same goes for response senders or for handlers - you can add a new file in /handlers and add your users, auth, products handlers etc. so that the structures remains flat, with maintaining clarity of what is where. 

Also, the responses that the server gives follow a standard, one for 200+ status codes, and another ones for the errors: 400+, 500+ status codes. This is implemented by the sender that you'll use to respond to requests.

// logs
// The logs directory path has been set to this root directory, but do not keep them here. Add your own path to a
// directory outside of the root project, like /var/log/myapp for Unix/Linux systems, where you can separate the concerns
// or create log rotation if needed.

## Setup

Make sure to first install the binaries that will generate the api docs and hot-reload the app.

```go install github.com/swaggo/swag/cmd/swag@latest```

and

```go install github.com/cosmtrek/air@latest```

Download the libs

```go mod download```

```go mod tidy```

Create a `.env` file in the root folder and use this template:
```
# DEV, STAGE, PROD
ENV=DEV
PORT=8080
VERSION=0.0.1

DB_HOST=localhost  #when running the app without docker
# DB_HOST=postgres # when running the app in docker
DB_USER=postgres
DB_NAME=postgres
DB_PASSWORD=postgres
DB_PORT=5432
```

If you start the app locally without docker make sure your Postgres DB is up.
Write `air` in terminal and the app will be up and running listening on whatever port you set in .env.

## Template Tour


## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

