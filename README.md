# Simple Messenger API 

## Overview

This is simple API that supports a web-application to enable two users to send short text messages to each other.

High-level components include:
- An API server written in Go, leveraging its native concurrency model (via goroutines & channels) and [gorilla/websocket](https://github.com/gorilla/websocket) package to handle server-side sockets. 
- MongoDB is used to persist chat histories.
- Docker is used to build & deploy.

## How to run the API server

If using Docker, run:
```bash
$ docker-compose up --build
```
\* requires [Docker Engine](https://docs.docker.com/engine/) and [Docker Compose](https://docs.docker.com/compose/install/).

If in dev mode, run:
```bash
$ go mod download
$ go run ./
```
\* requires [Go 15+](https://go.dev/dl/) & a local [MongoDB](https://www.mongodb.com/) instance running on default port.

## How to consume the API server

The API provides the following endpoints:

### Join a chat

```
GET (Upgrade: websocket) /ws 
```

**Code samples**

e.g. via [wscat](https://www.npmjs.com/package/wscat) in the terminal:
```bash
// terminal 1
$ wscat -c ws://localhost:8080/ws
> {"to": "pei", "from": "ave", "text": "hey mommy!"}
> {"to": "pei", "from": "ave", "text": "i am hungry, feed me!"}
> {"to": "pei", "from": "ave", "text": "and i want to play fetch!"}
> {"to": "pei", "from": "ave", "text": "but also please come snuggle me!"}
// terminal 2
$ wscat -c ws://localhost:8080/ws
> {"to": "ave", "from": "pei", "text": "hold on miss ave, i gotta finish this coding assignment."} 
```

e.g. via JS
```javascript
let websocket = new WebSocket("ws://" + window.location.host + "/ws");
websocket.send(
    JSON.stringify({
        from: "ave",
        to: "pei",
        text: "hey mommy!",
    })
);
```

### Retrieve recent messages
List the messages between a sender and a recipient in the last 30 days, in reverse chronological order.

```
GET /messages
```

**Parameters**

| Name | Type | In  | Description |
|---------|--------|---------| --------------------|
| to | string | query | The username of the sender |
| from | string | query | The username of the recipient  |
| limit | int | query | Number of messages returned. Default: 100 |


**Code samples**

e.g. via curl
```bash
// get last 10 messages from pei to ave:
$ curl http://localhost:8080/messages\?to\=pei\&from\=ave\&limit\=10 | jq '.[].text'

// get last 5 messages sent to pei:
$ curl http://localhost:8080/messages\?to\=pei&limit\=5 | jq '.[].text'

// get last 100 messages
$ curl http://localhost:8080/messages | jq '.[].text'
```

## Notes

Given the characteristics of a messaging API, I used WebSocket protocol to allow a two-way, persistent communication b/t the server and the client. 

Re: data modeling, the `Client` holds a websocket connection and a single instance of the `Server` type; the `Server` maintains a set of joined clients and broacasts messages to all clients. 
- `Server` has separate channels for joining clients, evicting clients, and broadcasting messages;
- `Client` has a buffered channel `buffer` for messages. One of its goroutines reads messages from this channel and writes messages to the websocket, while the other reads messages from the websocket and sends to the `Server`.

Decided to add a persistance layer using MongoDB.
- Some of my assumptions include: 1) write/read ratio about 1:1 for 1-1 chat; 2) only recent data is accessed & needs reasonable performance for random reads. 
- Currently has index on `timestamp` but should have a compound index on `_id` and `timestamp` 

Considerations/Limitations:
- supports 1-1 chat only, but the data models can be extended to support group chats later.
- ignores authorization, authenticaion and registration, or other security measures.
- no unit tests written due to time constraint.
- does not consider scalability
