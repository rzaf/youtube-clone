# youtube-clone

media storage and sharing api written in go using microservice architecture.

- for simplicity and data integrity used a central database for user related data except media files
- written in a request-response architecture
- services communicate to each other using gRPC

## microservices

### auth service

authentication microservice used for registering user and retrieving user JWT token and refreshing it

### gateway service

main api of app for user related data (users,medias,comments,playlists,likes,tags,followings)

### file service

- has a Rest api for recieving and uploading media files (photo,video,music)
- stores all media files and its meta data (other services communicate with file service to know media file url exist and belong to which user)
- media files are stored in system files and meta data is stored in a mongodb database
- handles video,music formating and compression (using ffmpeg)

### database service

- central database (postgres) for storing all user user related data (users,medias,comments,playlists,likes,tags,followings)
- other services communicate with database service to get user related data

### notification service

- recieves notifications from other services to send to user
- send emails to users for verifications and notifications

## api endpoints

protected routes require authentication with bearer tokens

### auth service routes

| type |        url         |protected|
|------|--------------------|-------|
| GET  |`api/login`         |&cross;|
| GET  |`api/register`      |&cross;|
| GET  |`api/refresh`       |&cross;|

### gateway service routes

#### USERS

| type |                 url                |protected|
|------|------------------------------------|-------|
| GET  |`api/users/{username}`              |&cross;|
| GET  |`api/users`                         |&cross;|
| GET  |`api/users/search/{term}`           |&cross;|
| POST |`api/users/{username}/verify/{code}`|&cross;|
| POST |`api/users/resend-email`            |&check;|
| PUT  |`api/users/{username}/profile-photo`|&check;|
| GET  |`/users/{username}/followings`      |&check;|
| PUT  |`api/users/{username}/channel-photo`|&check;|
| PUT  |`api/users/{username}`              |&check;|
|DELETE|`api/users/{username}`              |&check;|
| POST |`api/follows/{username}`            |&check;|
|DELETE|`api/follows/{username}`            |&check;|

#### MEDIAS

| type |                 url                      |protected|
|------|------------------------------------------|-------|
| GET  |`api/medias`                              |&cross;|
| GET  |`api/medias/search/{term}`                |&cross;|
| GET  |`api/medias/{url}`                        |&cross;|
| POST |`api/medias`                              |&check;|
| PUT  |`api/medias/{url}`                        |&check;|
|DELETE|`api/medias/{url}`                        |&check;|
| POST |`api/medias/{url}/tag/{name}`             |&check;|
|DELETE|`api/medias/{url}/tag/{name}`             |&check;|
| POST |`api/medias/{url}/playlists/{playlistUrl}`|&check;|
| PUT  |`api/medias/{url}/playlists/{playlistUrl}`|&check;|
|DELETE|`api/medias/{url}/playlists/{playlistUrl}`|&check;|

#### COMMENTS

| type |                 url                      |protected|
|------|------------------------------------------|-------|
| GET  |`api/comments/{commentUrl}`               |&cross;|
| GET  |`api/comments/medias/{url}`               |&cross;|
| GET  |`/comments/{commentUrl}/replies`          |&cross;|
| POST |`api/comments/medias/{url}`               |&check;|
| PUT  |`api/comments/{commentUrl}`               |&check;|
|DELETE|`api/comments/{commentUrl}`               |&check;|

#### LIKES

| type |                 url                      |protected|
|------|------------------------------------------|-------|
| POST |`api/medias/{url}/likes`                  |&check;|
|DELETE|`api/medias/{url}/likes`                  |&check;|
| POST |`api/comments/{url}/likes`                |&check;|
|DELETE|`api/comments/{url}/likes`                |&check;|

#### PLAYLISTS

| type |                 url                      |protected|
|------|------------------------------------------|-------|
| GET  |`api/playlists`                           |&cross;|
| GET  |`api/playlists/search/{term}`             |&cross;|
| GET  |`api/playlists/{url}`                     |&cross;|
| GET  |`api/playlists/{url}/medias`              |&cross;|
| POST |`api/playlists`                           |&check;|
| PUT  |`api/playlists/{url}`                     |&check;|
|DELETE|`api/playlists/{url}`                     |&check;|

### file service routes

| type |        url         |protected|
|------|--------------------|-------|
| GET  |`api/photos/{url}`  |&cross;|
| GET  |`api/videos/{url}`  |&cross;|
| GET  |`api/musics/{url}`  |&cross;|
| POST |`api/photos/upload` |&check;|
| POST |`api/videos/upload` |&check;|
| POST |`api/musics/upload` |&check;|

## Prerequisites

- ***Go*** (1.20)
- ***docker*** and ***docker compose***
- Protobuf compiler (`protoc`) & go plugins (optional):
  - install protocol buffer compiler [link](https://grpc.io/docs/protoc-installation/)
  - install protoc-gen-go and and protoc-gen-go-grpc by running `go install google.golang.org/protobuf/cmd/protoc-gen-go` and `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc`
  - have protoc and GOPATH in your PATH env

## build

### docker

1. Build Go Files and Docker Containers

   To build the Go binaries and Docker containers, use:

   ```bash
   make build 
   ```

   If protoc and the Go plugins are installed, you can generate Protocol Buffer files and build Swagger docs files along with the Go binaries:

   ```bash
   make all 
   ```

2. Create .env file

   ```bash
   cp .env.example .env 
   ```

3. Run Docker containers

   Start the Docker containers using:

   ```bash
   make run
   ```

   or simply run `docker compose up`

### docker swarm

Deploy using Docker Swarm

```bash
make swarm
```

## features

- video,music,photo storage and sharing
- video and music streaming
- swagger documentations (at `/docs`)
- user subscribing (following)
- comment and reply on medias
- media and comment likes
- searching users,medias,playlists
- multi media playlists
- email verification
- user upload limit
- notifications:
  - new media from following user (subscribing)
  - new comment on users media
  - new reply on users commnet
  - new like on users media,comment,reply
  - new follower (subscriber)

## todo

- image compression
- multiple video quality and bitrate
- video subtitle
- video thumbnail from video
- adding image to musics
