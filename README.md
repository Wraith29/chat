# Go Chat App

## Usage

**Server**

```sh
> go run main.go server
> # OR
> ./server
```

**Client**

```sh
> go run main.go client <client_name>
> # OR
> ./client <client_name>
```

## Features

**General**

- [ ] Graceful logging (with the app taking over the page, need to log to a file to ensure we can view logs in case of a crash etc)

**Server**

- [x] Allow users to connect
- [ ] Allow users to disconnect at will
- [ ] UI
  - [ ] See active connections
  - [ ] View Messages sent by a specific user
  - [ ] Force Disconnect user

**Client**

- [x] Connect to server
- [ ] Disconnect gracefully (`/commands`?)
- [ ] UI
  - [x] Send messages
  - [ ] View Messages (Scroll once reaching a certain amount?)

