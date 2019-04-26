# PrQL: SQL over HTTP
_Note: This project is still in development. It is __not safe__ for production environments._

[![CircleCI](https://circleci.com/gh/PrQL/prql/tree/master.svg?style=svg)](https://circleci.com/gh/PrQL/prql/tree/master)


#### What is PrQL?

PrQL is an ultra-lightweight layer that you place between your serverless application and your SQL database(s). 
What it does is allow a frontend application to gain access to your data without having to build out a backend.
Sometimes our needs are very simple and there isn't any purpose to building, maintaining, and standing up a REST API
to funnel data between the SQL database and the client. That's where PrQL steps in and allows you to securely query
your database from any client (given that client has a token). In addition to ease of access, it allows the client to
decide what data it will receive and how it will receive it since you pass along a plain SQL query.

#### What about GraphQL and JSON API?

You could consider PrQL to be a loose competitor to those two methods, but PrQL has different intentions. Although they 
both allow you to define relationships within the data, neither give as much precision as actual SQL. This is especially
problematic if you work in the realm of GIS. Certain databases/database extensions allow you perform  complex operations on 
the dataset before handing the data over to your application. Implementing an array of custom, complex operations for your 
app will most likely break REST conventions once you are all done writing your API. Why attempt to adhere to an ill-fitted 
specification when you know that, in your work domain, you will always be breaking that spec?

GraphQL and JSON API are great in the proper domain, but: JSON API is fancy REST. GraphQL is too simple yet too complex at 
the same time.

#### What about PostgREST?

PostgREST is brilliant, but again, it serves a different purpose than PrQL. In addition to different solution sets, PostgREST
is specific to PostgreSQL while PrQL is database agnostic. As long as you have a SQL database, PrQL will have a driver to 
connect to it. 


## Usage

- Define database locations
- Generate a token to authenticate against a defined database
- Add token to request headers in your client
- Query your service for data


### Setting up a database

First, you have to inform PrQL of your database's location. You save this information by giving that database a name.

```sh
sudo prql databases new \
            --name localpg \
            --driver postgres \ # possible options are postgres and mysql (will add more post-alpha)
            --port 5432
```

Then, to view the newly added database:
```sh
sudo prql databases list
```


### Generating an app token

Now you can generate a token that maps to a user in a database that you have previously defined. That way you have a 
secure way of allowing the client to let PrQL know which database server and user to execute queries on.

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

### Query your database

To tell PrQL which credentials to use, we pass a token as a query parameter. 

```sh
wget -nv -O- "localhost:1999?token=f04e79dc8d1dd5453da438366c6162fb&query=SELECT id, name FROM users WHERE login_attempts > 3"
```


## Setup

### Building

There are currently a two different ways of implementing the PrQL alpha:

#### Go

##### Dependencies
- [dep](https://github.com/golang/dep#installation)
- [staticcheck](https://github.com/dominikh/go-tools/tree/master/cmd/staticcheck)

```sh
mkdir -p "$GOPATH/src/github.com/prql"
cd "$GOPATH/src/github.com/prql"
git clone https://github.com/prql/prql && cd prql
make
```

#### Docker

```sh
make with-docker ARCH=$YOUR_DISTRO
```
Where `$YOUR_DISTRO` is one of the following: `darwin/amd64` `darwin/386` `freebsd/amd64` `freebsd/386` `linux/arm` `linux/arm64` `linux/amd64` `linux/386` `solaris/amd64`

You will find your binaries in the _build/_ folder with your distro name appended.

### Installation

On *nix systems, run `sudo ./install.sh`. This is a temporary measure for the alpha release. It sets up the necessary
system directories and files, installs binaries, and attempts to enable `prqld` as a systemd service.

On Windows systems, :shrugs: I haven't figured out what that needs yet. It's in alpha afterall.

### Running

If systemd exists on your system, `reboot` or `sudo systemctl start prqld`.

Otherwise, you will have to manage prqld on your own by executing `prqld &`.

To check if prqld is up and running, you can 
```sh
lsof -iTCP -sTCP:LISTEN | grep prqld
```
