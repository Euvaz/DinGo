# DinGo

DinGo is a Discord chat bot intended to provide various useful features.

## Overview

- [Usage](https://github.com/euvaz/dingo#usage)
- [Installation](https://github.com/euvaz/dingo#installation)
  - [Configuration](https://github.com/euvaz/dingo#configuration)
  - [Execution](https://github.com/euvaz/dingo#execution)
  - [Inline](https://github.com/euvaz/dingo#inline)
  - [Binary](https://github.com/euvaz/dingo#binary)
- [License](https://github.com/euvaz/dingo#license)

## Usage
This project assumes that Go is both installed, and added to your system's PATH.

## Installation

### Configuration

Firstly, it is necessary to update the environment file.

```sh
$ cp .env_sample .env
```

After copying over the .env sample file, set the values accordingly.

### Execution

There are two methods of running the bot, [Docker](https://https://github.com/euvaz/dingo#docker) (Recommended) or [Binary](https://github.com/euvaz/dingo#manual). Utilizing the docker method is recommended, as the docker-compose will automatically create a self-hosted [MariaDB](https://mariadb.org/) instance.

### Docker

When using this method, be sure to set the following values within `.env`:

```
...
PSQL_HOST=127.0.0.1
...
```

Build image and run as daemon:

```sh
$ docker compose up -d
```

If any changes are made to the environment file, the docker compose will need to be stopped, rebuilt, and started again:

```sh
$ docker compose down
$ docker compose build
$ docker compose up -d
```

### Manual

There are two methods of running the project.

1. Inline - Useful for debugging:

    ```sh
    $ go run main.go`
    ```

2. Binary - Recommended for faster execution times:

    ```sh
    $ go build main.go
    $ ./main.go
    ```

This will result in a `main` executable file being created, with the filetype dependent upon the system which the executable was generated.

## License

As of 18-Jan-2022, DinGo is fully open to the public, licensed under GPLv2.
