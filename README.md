## forum (backend)

### Usage

The REST API (backend part) of this application can be accessed at: [https://backend-13af.up.railway.app](https://backend-13af.up.railway.app)

### Tech Stack:

Python, Django, Django DRF, PyJWT, SQLite, Docker.

## Services

#### Users

- **POST /user/signup**: Registration.
  - Request Body:
    ```json
    {
        "username": "john_doe",
        "email": "johndoe@example.com",
        "password": "example_password"
    }
    ```
    
- **POST /user/login**: Authorization.
  - Request Body:
    ```json
    {
        "username": "john_doe",
        "password": "example_password"
    }
    ```
  - Response Body:
    ```json
    {
        "access_token": "example_token",
        "type": "Bearer"
    }
    ```
    
#### Posts

- **POST /posts**: Create.
  - Request Body:
    ```json
    {
        "title": "example_title",
        "content": "example_content",
        "categories": ["1","2","5"]
    }
    ```
    
- **PUT /posts/{post_id}**: Update.
  - Request Body:
    ```json
    {
        "title": "example_title",
        "content": "example_content"
    }
    ```
    
#### Comments

- **POST /posts/{post_id}/comments**: Create.
  - Request Body:
    ```json
    {
        "content": "example_content"
    }
    ```
    
- **PUT /posts/{post_id}/comments/{comment_id}**: Update.
  - Request Body:
    ```json
    {
        "content": "example_content"
    }
    ```
    
#### Reactions

- **POST /posts/{post_id}/react**: React to post.
  - Request Body:
    ```json
    {
        "is_like": 1
    }
    ```
    
- **POST /posts/{post_id}/comments/{comment_id}/react**: React to comment.
  - Request Body:
    ```json
    {
        "is_like": 0
    }
    ```
