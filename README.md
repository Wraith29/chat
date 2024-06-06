# Go Chat App

## Usage

**Server**

```sh
> go run cmd/server/server.go
> # OR
> ./server
```

**Client**

```sh
> go run cmd/client/client.go <client_name>``
> # OR
> ./client <client_name>
```

## Features

**General**

- [ ] Graceful logging (with the app taking over the page, need to log to a file to ensure we can view logs in case of a crash etc)

**Server**

- [x] Allow users to connect
- [x] UI
  - [x] See active connections
  - [x] View Messages sent by a specific user
  - [x] Force Disconnect user

**Client**

- [x] Connect to server
- [ ] Disconnect gracefully (`/commands`?)
- [ ] UI
  - [x] Send messages
  - [ ] View Messages (Scroll once reaching a certain amount?)

