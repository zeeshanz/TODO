# TODO

This is a TODO app as a part of the onboarding process which uses Golang, Go Fiber and PostgreSQL.

PostgresSQL is running in a Docker instance. The Go server is running on port 3001.

Make sure the directory which contains the files has correct persmissions and ownership.

Run `docker compose up -d` to start Docker.

Run `go run migrate/migrate` to create tables in PostgreSQL

Run the app using `go run main.go routes.go`

In a browser go to `localhost:3001` to run the app.
