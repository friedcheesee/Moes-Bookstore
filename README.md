[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/LECuYE4o)
# Moe's Bookstore

This is a fully functional bookstore written in golang, which uses PostgreSQL as a database. The app has many API endpoints which the users can call to access many functions like login, register, logout, search books, add to cart, remove from cart, buy books, download books from inventory. The app also has separate API's for admins to manage users. A Sample database has been provided in the 'sql-scripts' directory.(automatically configured)
The app opens on port 8080, and the database opens on port 5432, and returns responses in form of JSON so that it can be easily integrated with a webapp.

## Index:
1. [Features](#features)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Usage](#usage)
4. [Working](#working)
5. [Description](#description)
6. [API Endpoints](#api-endpoints)



## Features

- Uses parameterised SQL queries in code to prevent SQL injections
- Uses transactions to maintain the ACID property of the database.
- Has appropriate error handing methods in every function
- The app logs everything in detail inside the app.log file, including events and errors.
- Multiple clients can connect to the webserver at once, because the app uses cookies to identify each user.
- Uses Nginx as reverse proxy to serve the clients.
- Provides recommdation after buying a book.
- The app has been containerized, along with the database.
- App uses bcrypt and a random 32bit key to hash the passwords stored in the database.

![working](/assets/working.png)


## Prerequisites

- Golang
- Docker (if running in docker)
- Postman (Free software for calling API)
- PostgreSQL
- dependencies (included in go.mod)

The project was developed using psql (PostgreSQL) 15.4, and golang 1.21.0.

The dockerfile contains all the commands to make an image and the docker-compose.yml is included to help you run it in docker.

If required, the docker images are also available on docker hub: 
[Moe store](https://hub.docker.com/repository/docker/friedcheese/jao-moe)
[postgres](https://hub.docker.com/repository/docker/friedcheese/postgres)


## Installation

To import a Postman configuration,
1. Open Postman
2. Click 'Import' on the top left corner
3. Click 'select files'
4. Browse to the postman_collection provided in the root directory.


With Docker:

1. Clone this repository
2. Open a terminal in the root folder
3. run the command:

```bash
 'docker compose build'
 ```

4. next run the command, 

```bash
'docker compose up -d'
```
5. Open Postman and import the configuration provided.
6. Start sending requests via Postman, starting from Login/Register.

[A short video guide to install with docker](https://youtu.be/G8YjCzNMhLY) 
   
Without docker:

1. Clone this repository
2. Open a terminal in the root folder
3. Run the script file 'init.sql' included in 'sql-scripts' to setup relations and some pre generated data, steps to run a sql file:

```bash
\i 'D:/Project/moe/sql-scripts/init.sql'
```
where D:\Project\moe\sql-scripts\init.sql is the path to the script file. In case of errors, the .sql can be opened in notepad, and all the commands required to make the relations are included. Make note of the backward slashes in the command.

4. Open 'loginauth.go' to setup your database connection variables, such as username, password, database name, host, port etc. in Adminconnect() function.

5. Open 'main.go' file
6. Comment the line "db = ah.Newconnect()"[Line 41] and uncomment the line "db=ah.Adminconnect()"[Line 40]. We need to use Adminconnect() when connecting to the local database, Newconnect() will be used when connecting to the docker container(it is already setup).
7. run the command:
```bash
'go mod download'
``` 
8. Start the application by the command
 ```bash
'go run .'
```
(Allow the program in the firewall, if asked)

9. Now you must be able to see 'Connected' in the terminal, which means the database connection has been made.
10. Open Postman and import the configuration provided.
11. Start sending requests via Postman, starting from 'login'/register.


## Usage

- This project should work perfectly fine in docker.
Most of the endpoints are not accessible without logging in to the database. 

- Login credentials for an admin: email:fried@mail.com password: 'abcd'.

- Login credentials for a normal user: email:apple@mail.com password : 'abcd' or you can register a new user of your choice.

- All the passwords in the current data are 'abcd' for simplicity of testing.

- Sample user path - 

1.Login
2.Search Books by changing params in postman:
![parameters](/assets/parameters.png)

and click 'Send' on the right.

3. Add to cart using book ID from results:
![body](/assets/body.png)

4. Buy books by directly sending the request.

5. View inventory - Has the link to download any of your owned books.

6. Logout

The app has many more functions, all are predefined with body and parameters in the collection provided.



## Working

1. The app starts by initialising the log files, reading the environment variables, starting the cookiestore and creating a connection to the database.
2. The app then starts the chi router and defines all the endpoints.
3. The user can then send requests to the endpoints via Postman.
4. The app will then call the respective handler function, which will call the respective function to perform the required operation on the database.
5. The function will then return the data to the handler, which will then return the data to the user in the form of JSON.
6. The app will log all the events and errors in the app.log file.
7. User endpoints go through a middleware wrapper which checks if the user is logged in or not, and if the user is logged in, it will call the respective handler function, else it will return an error.
8. Admin endpoints go through a middleware wrapper which checks if the user is logged in and is an admin or not, and if the user is logged in and is an admin, it will call the respective handler function, else it will return an error.
9. Deleted accounts are marked as 'false' in the 'active' column of the users relation, hence their personal data cannot be fetched by any function, but the data is still present in the database, such as reviews.


## Description

.env: This file contains all the environment variables required for the app to run, such as database credentials, port number, session key etc. Change the values in this file to suit your needs, the variables inside the current env file are:
    
 ```bash
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=moe
COOKIE_KEY=65f3274ea398ef20fbe28512ac992a1f4eac9cae87da854d1f135f4f18e4e916
```
- GitHub hides .env files by default, but i have created the .env file manually in the repository.

- main.go - Initialises log files, reads variables from environment file, starts the cookiestore to store sessions and creates a connection to the database. Contains all the routes and endpoints for API calls, uses chi router.

- handlers.go - Has middleware wrapper functions and handlers, which will call the respective functions when it receives a request on their respective endpoints, as defined in main.go

- books.go - Contains functions which majorly deal with operations on the database related to books, such as adding to cart, removing, getting recomendations etc

- loginauth.go- Contains functions required for authorization, such as fetching ID from the database, hashing passwords using bcrrpt,  validating email, storing credentials for a user.

- admin.go - groups all related functions to admin operations.

- log.go - This has a package clled 'moelog' which is imported by other packages to keep track of events and errors throughout the program.

- models.go - contains structures to help handlers/functions return data in a formatted manner.

## API-Endpoints

Some API definitions:

### Login

- **Endpoint:** `/login`
- **Method:** `POST`
- **Description:** Authenticate a user by providing their email and password.
- **Request Parameters:** `email`, `password`
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/login \
    -d "email=fried@mail.com" \
    -d "password=abcd"
```
### Logout

- **Endpoint:** `/user/logout`
- **Method:** `POST`
- **Description:** Logs user out
- **Request Parameters:** None
- **Example Request:**
```bash
    curl -X POST http://localhost:8080/user/logout
```

### Register

- **Endpoint:** `/reguser`
- **Method:** `POST`
- **Description:** Register a new user by providing their email, password, and username.
- **Request Parameters:** `email`, `password`, `username`
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/reguser \
  -d "email=newuser@mail.com" \
  -d "password=abcd" \
  -d "username=its me"
```
### Inventory

- **Endpoint:** `/user/inventory`
- **Method:** `POST`
- **Description:** Get the inventory of books for the authenticated user.
- **Request Parameters:** None
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/user/inventory
```
### Search Books
- **Endpoint:** `/user/search`
- **Method:** `GET`
- **Description:** Search for books by providing a search query.
- **Request Parameters:** `query`, `genre`, `author`
- **Example Request:**
```bash
  curl -X GET http://localhost:8080/user/search?query=&author=&genre=drama
```
### Add to Cart
- **Endpoint:** `/user/cart/add`
- **Method:** `POST`
- **Description:** Add a book to the user's cart.
- **Request Parameters:** `bookid`
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/user/cart/add \
  -d "bookid=38"
```
### View Cart
- **Endpoint:** `/user/cart/view`
- **Method:** `POST`
- **Description:** View the user's cart.
- **Request Parameters:** None
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/user/cart/view
```
### Delete from Cart
- **Endpoint:** `/user/cart/delete`
- **Method:** `POST`
- **Description:** Delete a book from the user's cart.
- **Request Parameters:** `bookid`
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/user/cart/delete \
  -d "bookID=15"
```
### Buy Books
- **Endpoint:** `/user/buy`
- **Method:** `POST`
- **Description:** Buy all the books in the user's cart.
- **Request Parameters:** None
- **Example Request:**
```bash
  curl -X POST http://localhost:8080/user/buy
```