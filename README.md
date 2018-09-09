# Magpie

An agent that collects configuration information from running Docker containers.

## How it works

Magpie uses the Docker event API to retrieve the container's environment variables
and stores them in a PostgreSQL database along with the service name, version, and the
container's created timestamp.

It uses a blacklist to filter out variables
and a whitelist to include the variable in clear text.
If a variable is not in either list it will be masked.

A Docker container's variables are only saved if there aren't already
a set of variables with the same service name, version, and timestamp.

## Configuration

Use a configuration file to set database credentials and so.
If no file is specified, magpie will use the Postgres defaults
unless `~/.magpie.conf` exists.

Example configuration:

```
username=yourusername
password=secretpass
host=database-hostname
port=5432
database=mydb
```

## Database migration

Use `magpie -migrate` to migrate the database.

If a migration fails and the program complains that the migrations
are dirty, use `-force <int>` (where `int` is the ID of the failed migration)
to force the migrations to a good state.

For more on how to write migrations, [look here](migrate/steps/README.md).

## Running

Magpie accepts the following flags:

`-config string`

Specify an alternative path to the configuraion file
(default is `~/.magpie.conf`).

`-init`

Magpie will scrape all currently running Docker containers
when it starts before listening for events.

`-migrate`

Migrate the database.

`-migrate -force <int>`

Force the migration schema to the given version.

## Development & hacking

To build locally, you need the `go-bindata` tool installed.
This is used for bundling database migrations in the binary.

You can install it by running

`go get -u github.com/jteeuwen/go-bindata/...`

Then run `make migrate` to generate the binding file.

[More information on `go-bindata`](https://github.com/golang-migrate/migrate/tree/master/source/go_bindata).

Using the `Makefile` to build should suffice;
`make build` compiles the binary
and `make install` installs it on your system.
