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

```bash
npm install react-native-xportal
```
or 
```bash
yarn add react-native-xportal
```

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
å```

## Description

as flat as possible
The library works as a react-native substitute for [mx-sdk-dapp](https://github.com/multiversx/mx-sdk-dapp/tree/main), it helps mobile apps connect and interact with the XPortal wallet, including providing the necessary account information (balance, tokens, address etc) and signing transactions, messages or custom requests, thus abstracing all the processes of interacting with users' wallets. On connect, sign transaction or other actions, XPortal app will automatically be opened through deeplinking to complete the intended action. 

The library has 2 main modules: `core` and `UI`. The `core` modules gives you the functions to connect, sign transactions and other so that you can call them anywhere you need. The `UI` modules exports buttons for ease of use, they also use the `core` module under the hood.

***Note: This library has all the basic functionalities for interacting with the XPortal Wallet. New functionalities can be added - if you want to contribute, please see the [Contributing](#contributing) section.***

## Usage
The library needs to be initalized first in order to work, see example below.
```typescript
import { XPortal } from 'react-native-xportal';

const callbacks = {
      onClientLogin: async () => {
            console.log('on login');
      },
      onClientLogout: async () => {
            console.log('on logout');
      },
      onClientEvent: async (event: any) => {
            console.log('event -> ', event);
      },
};

try {
      await XPortal.initialize({
            chainId: 'd',
            projectId: '<wallet connect project ID>',
            metadata: {
                  description: 'Connect with X',
                  url: '<your website>',
                  icons: ['<https://img.com/linkToIcon.png>'],
                  name: '<name>',
            },
            callbacks,
      });
} catch (error) {
      console.log(error);
}
```
You need to  have a WalletConnect project ID. To get one see: https://cloud.walletconnect.com/app. Also, make sure to have valid data in your metadata key, otherwise the XPortal wallet will show a "Unexpected Error" when redirecting to it for login.

### Core
#### Login
```typescript
import { XPortal } from 'react-native-xportal';

try {
      await XPortal.login();
} catch (error: any) {
      throw new Error(error.message);
}
```
This will connect your app to user's XPortal app and his account.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

