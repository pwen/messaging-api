# Simple Messenger API 

## Overview

This is simple API that supports a web-application to enable two users to send short text messages to each other.

High-level components include:
- An API server written in Go, leveraging its native concurrency model (via goroutines & channels) and [gorilla/websocket](https://github.com/gorilla/websocket) package to handle server-side sockets. 
- MongoDB is used to persist chat histories.
- Docker is used to build & deploy.

## Notes

Given the characteristics of a messaging API, I used WebSocket protocol to allow a two-way, persistent communication b/t the server and the client. 

Re: data modeling, the `Client` holds a websocket connection and a single instance of the `Server` type; the `Server` maintains a set of joined clients and broacasts messages to all clients. 
- `Server` has separate channels for joining clients, evicting clients, and broacasting messages;
- `Client` has a buffered channel `send` for messages. One of its goroutines reads messages from this channel and writes messages to the websocket, while the other reads messages from the websocket and sends to the `Server`.

Decided to add a persistance layer using MongoDB.
- Some of my assumptions include: 1) write/read ratio about 1:1 for 1-1 chat; 2) only recent data is accessed & needs reasonable performance for random reads. 
- Currently has index on `timestamp` but should have a compound index on `_id` and `timestamp` 

Considerations/Limitations:
- supports 1-1 chat only, but the data models can be extended to support group chats later.
- ignores authorization, authenticaion and registration, or other security measures.
- no unit tests written due to time constraint.
- does not consider scalability
