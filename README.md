# fairmoney-test

### Prerequisites

1. **Go 1.20** or **lastest version** already installed on your local machine.
2. mongodb server

### Run Applications from Root

1. Ensure you mongodb instances are running
2. Create and populate `.env` files on each applications root with its keys and corresponding values as listed in `sample.env`
3. Run from each applications root directory

```bash
$ go run main.go
```

### Run Project as Docker container

1. Ensure you mongodb instances are running
2. Create and populate `.env` files on each applications root with its keys and corresponding values as listed in `sample.env`
3. Build Docker Image by running this command on the project root

```bash
$ docker compose up --build
```

### Testing

1. Automated unit and integration tests done with golang's builtin [`testing`](https://pkg.go.dev/testing) package.

To run all tests:
cd into each application's root folder and run

```bash
$ go test -v  ./tests/...
```

### Documentation

each of the applications have open api documention on this route

{baseUrl}/swagger/index.html

### Authentication

The app is currently not authenticated, but I would have used a middleware to implement authentation using bearer token header and jwt token.
