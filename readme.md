export ENV=DEV
export VERSION=0.0.1
export PORT=3001



// # errors 
// # maybe switch to sqlx?
// # contexts and tracing ?
// # logs ?
// # onshutdown
// # swagger
// # migrations
// # switch from dev to prod
// --->for me: to switch to prod, when you put the app on the production server, add the variables from env to terminal,
// or to a bash profile, this way the app will have those variables in the terminal it is using to run
// and will be accessible to it with os.Getenv. For dev, .env is enough and is good for sensitive info, yaml can be used
// but it's more for things that are hierachical and does not contain sensitive information.

// dockerize
// users

// what is .golangci.yaml ?

// I did not add the ability to run multiple servers since it would be too much for a template that tries to be clean and clear
// but if you need you can surely add it without modifying much
// or other specific features
// this is supposed to give you a solid base from which you can add your specific features

// logs
// The logs directory path has been set to this root directory, but do not keep them here. Add your own path to a
// directory outside of the root project, like /var/log/myapp for Unix/Linux systems, where you can separate the concerns
// or create log rotation if needed.

// I added two examples of outputting logs to file in /books and /book/{id}

// Set up for production is to add the variables in `~/.bash_profile` or `/etc/environment` or directly in terminal.

// go install github.com/swaggo/swag/cmd/swag@latest
// go install github.com/cosmtrek/air@latest