# Download & Build

## System requirements

The Ledger performance benchmarks run Ledger on 8 vCPU, 16GB RAM, 50GB SSD GCE instances, but any relatively modern
machine with low latency storage and a few gigabytes of memory should suffice for most use cases.

## Download the pre-built binary

The easiest way to get Ledger is to use one of the pre-built release binaries which are available for OSX, Linux, and Windows.
Instructions for using these binaries are on the [GitHub releases page][github-release].

## Build the latest version

For those wanting to try the very latest version, build Ledger from the `master` branch. [Go](https://golang.org/) version 1.13+
and [Makefile](https://www.gnu.org/software/make/manual/make.html) are required to build the latest version of Ledger.
To ensure Ledger is built against well-tested libraries, Ledger vendors its dependencies for official release binaries.

To build `Ledger` from the `master` branch using the official `build` script:

```sh
$ git clone https://github.com/danielnegri/tokenapi-go.git
$ cd tokenapi-go
$ make build
```

To build `Ledger` from the `master` branch using the official Docker `build` script:

```sh
$ make docker-image
```

## Test the installation

Check the Ledger binary is built correctly by starting Ledger and setting all
environment variables.

### Configure environment variables

Ledger requires configuration from environment variables or command line parameters.
The [Makefile] script reads and injects the environment variables from a `.env` file before running any command.

[Makefile#L40](../Makefile#L40)
```
# Inject env file
include .env
export $(shell sed 's/=.*//' .env)
```

Make sure to have a `.env` file in the root folder before starting the Ledger. You can copy and edit the sample included.

```sh
$ cp .env.sample .env
# Make sure to edit the database address, username, password, and token URL.
$ open .env
```

### Starting Ledger

If Ledger is built without using `make build`, run the following:

```sh
$ ./bin/ledger serve
```
If Ledger is built with Docker using `make docker-release`, run the following:

```sh
$ docker-compose run
```

### Generating tokens

Run the following:

```sh
$ curl -s -XPOST http://localhost:8080/api/v1/tokens?size=10
ERR: Hv58x-yQTRF1R71oMlGjKQ
OK : ijkr2lXOkM1EElPSDQFkeg
OK : 7fUa0EQMI6ZwjNco2YUOGA
ERR: MBvV7OncnJUCt-JK03z7mw
OK : eW8l673PRTjmUotzOgMDSQ
OK : QOWoKH7vXS9YLo704n72Hg
OK : LQ3SjTITNWVBsGD7J1htHw
OK : vn4QXi6xEI4BUB8phZKm7g
OK : pNXBNMG6Zsdjw2Gfueh6Nw
OK : _angmchPgYOeFU54dw1C9Q
```

If tokens are printed, then Ledger is working!

### Heartbeat

```sh
$ curl -s http://localhost:8080/health/heartbeat
AdhereTech Ledger Service @ 2020-07-14T18:12:43Z
```

[github-release]: https://github.com/danielnegri/tokenapi-go/releases/
[go]: https://golang.org/doc/install
[build-script]: ../build
[cmd-directory]: ../cmd
[example-hardware-configurations]: op-guide/hardware.md#example-hardware-configurations
