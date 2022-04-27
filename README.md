<div align="center">
<br><div><img height="118" src="https://assets.conlang.dev/images/brand/conlang-dev-logo.svg" alt="conlang.dev"></div>

**conlang.dev** is a tool for conlangers and linguists to document their languages.
</div>

## 🐻 About
This project is inspired by a desire for more open-ended tools that allow creative freedom in defining the limits of languages.
Inspired by tools like [Conworkshop](https://conworkshop.com) and [Lexique Pro](http://lexiquepro.com/), **conlang.dev** aims to sit somewhere between the two, with an extra focus on free-flowing documentation.

This repo is for conlang.dev's API. It's written using [Go](https://go.dev/) and is intended to be used in conjunction with a [MariaDB](https://mariadb.org/) database.

## 🏄 Getting started
With Docker, it's easy to get started with a dev environment.
```sh
# Run this to start the services up
# On first run, this will compile the API from source.
docker-compose up -d

# Run this to take them down again.
docker-compose down

# Run this to see running logs of the containers.
docker-compose logs -f

# Run this to rebuild just the API if there have been changes.
docker-compose up -d --build api
```

Make sure you've filled out the required environment variables! You can copy the base from `.env.example`.
```sh
# Copy the example file to .env
cp .env.example .env

# Edit it in your favourite editor! (sorry emacs-lovers)
vim .env
```

### Without Docker
Don't want to use Docker? No problem. Make sure you have Go installed.
```sh
# Compile the project into an executable
go build -o conlangdev github.com/conlangdev/conlangdev/cmd/conlangdev

# Run it!
./conlangdev run
```

You'll need an instance of MariaDB (at least version 10.5) for the database.

You'll also need the following variables in your environment:
* `CONLANGDEV_JWT_SECRET` - a string secret used to generate JWTs for authentication
* `CONLANGDEV_ADDR` - the address to listen on (`host:port`) - this can be just a port e.g `:8000`
* `MARIADB_HOST`, `MARIADB_DATABASE`, `MARIADB_USER`, `MARIADB_PASSWORD` for connecting to the database.

## 🐶 Developing
Make sure you write a migration for any changes to modelling.

Try and stick to the incrementing number convention for the migration filenames.
```sh
# Create a blank new migration file
touch sql/migrations/000x_my_migration.sql
```

If you change what external modules are in use such as by introducing a new module to the codebase, make sure you re-generate the `go.mod` and `go.sum` files. It's important for building the application that these files are accurate.
```sh
# Generate go.mod and go.sum files
go mod tidy
```

## 📚 License
conlang.dev is licensed under the permissive [BSD 2-Clause](https://opensource.org/licenses/BSD-2-Clause) license. You can read it in full in the [license](LICENSE) file.