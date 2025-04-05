## forum

### Usage: how to run

Prior to this, you will need to have docker and docker-compose installed on your machine.

- Open the terminal/command line and clone the repository to your local machine.

```console
git clone https://github.com/itelman/forum-nfactorial.git
```

- Open the repository from terminal.

```console
cd forum-nfactorial
```

- Run the following command:

```console
docker-compose up --build
```

- You will be able to access the service at: http://localhost:8080

Furthermore, the app can be accessed at: [https://forum-13af.up.railway.app](https://forum-13af.up.railway.app). The REST API (backend part) is at: [https://backend-13af.up.railway.app](https://backend-13af.up.railway.app).

### Tech Stack:

- Backend (REST API): Python, Django, Django DRF, PyJWT, SQLite, Docker.
- UI/UX (user-friendly interface): Golang, JavaScipt, HTML, CSS, Docker.

Additionally, Docker-Compose is used to combine the two services into a single app.

## Objectives

This project consists in creating a web forum that allows:

- creating posts and comments.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

#### Authentication

In this segment the client is able to `register` as a new user on the forum, by inputting their credentials. A `login session` is created to access the forum and be able to add posts and comments.

JWT tokens are used to log users into the REST API, encrypted and then applied in the UI/UX interface as a cookie value. Each created user session with such token is valid for 24 hours.

Instructions for user registration:

- Asks for email
  - When the email is already taken returns an error response.
- Asks for username
- Asks for password
  - The password is hashed when stored

The forum is able to check if the email provided is present in the database and if all credentials are correct. It will check if the password is the same with the one provided and, if the password is not the same, it will return an error response.

#### Communication

In order for users to communicate between each other, they are able to create posts and comments.

- Only registered users are able to create posts and comments.
- When registered users are creating a post they can associate one or more categories to it.
- The posts and comments are visible to all users (registered or not).
- Non-registered users are able to see posts and comments (not like and dislike them).

#### Likes and Dislikes

Only registered users are able to like or dislike posts and comments.

The number of likes and dislikes are visible by all users (registered or not).

#### Filtering and User Activity

A filter mechanism is implemented, that allows users to filter the displayed posts by categories. Filtering by categories is dispalyed as subforums for each category. A subforum is a section of an online forum dedicated to a specific topic.

Additionally, logged users can view their recent activity on the website, such as:

- created posts
- liked/disliked posts
