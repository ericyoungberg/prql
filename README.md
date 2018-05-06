# PrQL: SQL over HTTP

_Note: This project is still in development. It is __not safe__ for production environments._

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
problematic if you work in the realm of GIS. Certain databases/database extensions allow you perform complex operations on 
the dataset before handing the data over to your application. Implementing an array of custom, complex operations for your 
app will most likely break REST conventions once you are all done writing your API. Why attempt to adhere to an ill-fitted 
specification when you know that, in your work domain, you will always be breaking that spec?

GraphQL and JSON API are great in the proper domain, but: JSON API is fancy REST. GraphQL is too simple yet too complex at 
the same time.


## Setup

### Building

There are currently a few different ways of implementing the PrQL alpha:

#### Download

There are binaries for most modern distros:

#### Go

If you have go on your system, you can `go get` this repo and run `make install`.

#### Docker

If you have Docker on your system, you can build for your specific system by running
```sh
make with-docker ARCH=$YOUR_DISTRO
```
Where `$YOUR_DISTRO` is one of the following: darwin/amd64 darwin/386 freebsd/amd64 freebsd/386 linux/arm linux/arm64 linux/amd64 linux/386 solaris/amd64 windows/amd64 windows/386
ou will find your build in the _build/_ folder with your distro name appended.

### Installation

On *nix systems, run `sudo ./install.sh`. This is temporary measure for the alpha release. What it does is setup the necessary
system directories and files, installs binaries, and attempts to enable `prqld` as a 

On Windows systems, :shrugs: I haven't figured what that needs yet. It's in alpha afterall.

### Running

If systemd exists on your system, `reboot` or `sudo systemctl start prqld`.
Otherwise, you will have to manage prqld on your own by executing `prqld &`.

To check if prqld is up and running, you can 
```sh
lsof -iTCP -sTCP:LISTEN | grep prqld
```

Now you're ready to follow the usage guide below!


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
