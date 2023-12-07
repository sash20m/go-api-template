# Golang Rest API Template
> Golang Rest API Template with clear, scalable structure that can sustain large APIs.

## Table of Contents

- [Features](#Features)
- [Directory Structure](#Directory-Structure)
- [Description](#Description)
- [Setup](#Setup)
- [Template Tour](#Template-Tour)
- [License](#license)

## Features

- Standard responses for success and fail requests
- Swagger API documentation
- Sqlx DB with Postgres - but can be changed as needed.
- Standard for custom errors
- Logger for console and external file.
- Migrations setup
- Hot Reload
- Docker setup
- Intuitive, clean and scalabe structure
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
- /migrations --> Migrations to set up the database schema in your db.
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

I have spent a while looking for Go Api templates to automate my workflow but I had some trouble with the ones I've found. Some of them were way too specific, sometimes allowing just one set of handlers or one model in db to exist by default, in which case I had to rewrite and restructure it to make it open for extension, to add more handlers, DBs etc. Others were way to complex, with a deep hierarchical structure that didn't make sense (to me at least), half of which I would delete afterwards. So I wanted something that is as flat as possible but still having some structure that is easily extendable when needed, and also has the basic functionality set up so that I only need to add to it. This template is my attempt at achieving that.

To keep it simple, the template creates a CRUD of books which then can be deleted for your specific handlers, but they show the way the api works generally. The server uses `gorilla/mux` as a router, `urfave/negroni` as a base middleware, `sirupsen/logrus` as a logger and `unrolled/render` as the functionality to format the responses in any way you want before automatically sending them. The API also uses `unrolled/secure` to improve the security. For database management `jmoiron/sqlx` is used to improve the ease of use.

The app uses Air as a hot-reloader and Docker if you need it. You can start it in both ways (see [Setup](#setup)).
For environment variables `joho/godotenv` has been used instead of a `.yaml` file for security considerations. Again, everything is extendable, so you can add a `.yaml` file if you need more hierarchical structure or your environment variables don't need to be secure.
If you decide to not use Docker, in dev mode you use the variables from .env, and in production you add them in the terminal or in `~/.bash_profile`/`/etc/environment`. If you do decide to use Docker, keep the variables in .env in development mode, and add another .env on your production server with the prod variables. The .env file is not version controlled so there will not be conflicts. I chose to go with this approach from hundreds because I tried to hit the middle - simple and relatively secure, in a way that *most* people will use this template. If the api becomes large and you have specific needs for your case, you can add the variables in prod in command-line, in docker volumes, in a secret manager, in your kubernetes/docker swarm or any other way you want/need. But for most people and most cases, this approach will be more than enough.

I made sure you can add to it and modify without any pain, so for example, you can add one more db without modifying anything from the existing code, and also you can change the current db (Postgres) with any other also by not modifying anything besides creating the connect function for your new db. Same goes for response senders or for handlers - you can add a new file in /handlers and add your users, auth, products handlers etc. so that the structures remains flat, with maintaining clarity of what is where. 

Also, the responses that the server gives follow a standard, one for 200+ status codes, and another one for the errors: 400+, 500+ status codes. This is implemented by the sender that you'll use to respond to requests.

## Setup

Make sure to first install the binaries that will generate the api docs and hot-reload the app.

```
go install github.com/swaggo/swag/cmd/swag@latest
```
and
```
go install github.com/cosmtrek/air@latest
```

Download the libs
```
go mod download
```
```
go mod tidy
```

Create an `.env` file in the root folder and use this template:
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

Don't forget to rename the module and all the imports in the code from my github link to yours.

## Template Tour
Going from top to bottom, in the `/cmd` folder you'll see the entrypoint of the server. In its main.go you'll see the app loading the env and making a config struct which is based on the `/config` also present the root folder. This config will then be passed to the server, it's the only config and you probably don't need another one, just add to it if you want to pass more info to the server from env. The goroutine at the end of the function has the purpose of running a `OnShutdown` function when the server crashes or panics. You can add anything you want in the `OnShutdown` function, I left it empty but with some comments. In the goroutine we make a instance of our server struct which includes our whole server with the db, handlers and everything packed together. Also in `/server` you have a `Makefile` which builds the app properly in the same folder.

The `/internal` is where the actual code lies. The server creation is in the `server.go` (the `AppServer` struct that includes everything I said earlier). In the `Run` function of this struct we have all the server setup, adding the config info received from `/cmd`, the creation of the db, the router, middlewares and handlers, running the migrations and at the end we start everything up. The `Sender` and `Storage` fields on the struct are injected from `handlers.Handlers` thus the handlers will have access to them. The other 3 functions attached to the struct are `OnShutdown` that I mentioned earlier, `NotFoundHandler` and`NotAllowedHandler` that are 2 special handlers that handle the 404 and 405 status codes.

In `/errors` you have the standard for custom errors. You use this struct when you desire to respond to a request with some specific error info. This struct will be passed to the `Sender` which we'll see in a bit, and there the `Sender` will structure the response based on a specific standard. Below this there are the sentinels errors that you'll use throughout your app. The `Sender` accepts both a `string` as an error, an normal struct that implements the `error` interface and an `Err` struct that implements the `error` interface but has additional keys in the struct.

In `/handlers` you have the `handlers.go` file which creates the handlers object and all the dependencies it needs, which then itself is injected into `AppServer` in `server.go` that I mentioned above. Now, the others `.go` files in this folder should be specific to each handler group you have. By default there's one for books in `books.go`, but if you need handlers for users, you create a `users.go`, for auth an `auth.go` and so on. All the functions there are received by the `Handlers` struct in `handlers.go`.

`/middlewares` implements custom middlewares that are added into `negroni` at the start of the server in `server.go`.

In `/model` you first have the `base_model.go` which is the base for the other models. And in the same `/handlers` spirit, each file is responsible for its part of the app, so `books.go` will have the database model for books, the one that reflects the books table in db, and below it you have different structs that define how responses and requests for each endpoint should look like (since you don't always want to send the whole model struct as a response and also sending it with zero values is no good). 

In `/storage` there's the `storage.go` which defines how the storage should look like, thus being able to switch databases. The new database should implement the `StorageInterface` and it's good to go. `postgres_db.go` create a instance of this storage with a postgres db in the `Storage` struct as seen in `server.go` when we created the db. For a new db just create another `my_db.go` and create a `Storage` instance with that particular db. In the `Storage` struct you can add a `redis` db and so on if you want more dbs. Again in the same spirit, `books.go` implements the `StorageInterface` specifically for books, and all the methods are received by the `Storage` struct that is eventually used in handlers. For users you'll have `users.go` and there will be all the db operations specifically on users and so on.

In `\pkg` you have the `httputils` and there's the `http_response.go` which creates the standard for responses. The `Sender` struct in injected in the handlers and will be used there whenever you send a response. I made this struct with the idea of automation and also with the idea of extension - by default there's the `JSON` response, but you can add below HTML format, XML or any other you need. `s.Render` has a bunch of them. The `JSON` function takes the `http.ResponseWrites`, the statusCode and the struct you want to send back to the user, and based on the format defined above in the file, sends the response. 
The `/loger` in pkg defines the logger for the app, you can use this logger to log to console or to the `/logs` file in the root folder, you can see in handlers how it is used. The logs directory path has been set to this root directory, but do not keep them here. Add your own path to a directory outside of the root project, like `/var/log/myapp` for Unix/Linux systems, where you can separate the concerns or create log rotation if needed.

Thus you have all the minimum functionality you always need to get started on an Api that can grow large with time.

The other remaining files are self-explanatory.

Happy coding.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

