# PrQL: SQL over HTTP

_Note: This project is still in development. It is __not safe__ for production environments._

## Setup


## Usage

- Define database locations
- Generate a token to authenticate against a defined database
- Add token to request headers in your client
- Query your service for data


### Setting up a database
```sh
sudo prql databases new \
            --name localpg \
            --driver postgresql \
            --port 5423
```
Then, to view the newly added database:
```sh
sudo prql databases list
```
### Generating an app token
```sh
sudo prql tokens new \
         --host localpg \
         --database myapp \
         --user myapp_user
```
To view the tokens you have created, you follow the same pattern as the command above:
```sh
sudo prql tokens list
```
