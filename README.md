## Welcome Note

```
Everyone is welcome for code contributions or code enhancemnets.
please give a star, fork the repo and create PR for your code contributions
```
## db commands

```
create database userdb;
use userdb;
create table users(
  user_id int not null auto_increment primary key,
  username varchar(45) not null,
  password varchar(100) not null,
  email varchar(45) not null
  );
```
## run the application with air live reloading

```
air
```
## run the application with go run

```
go run main.go
```
## signup request

```
curl --location 'localhost:8080/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "user1",
    "password": "password1",
    "email": "user1@gmail.com"  
}
'
```
## signin request

```
curl --location 'localhost:8080/signin' \
--header 'Content-Type: application/json' \
--data '{
    "username": "user1",
    "password": "password1"
}
'
``` 

## profile request

```
curl --location 'localhost:8080/profile' \
--header 'Cookie: session_token=7e0ce38b-f86b-48fe-92f1-f27341e4d67a'

# you can get session_token from signin request
```
## signout request

```
curl --location 'localhost:8080/signout'
```
