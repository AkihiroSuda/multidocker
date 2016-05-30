# Multidocker

[![Build Status](https://travis-ci.org/AkihiroSuda/multidocker.svg?branch=master)](https://travis-ci.org/AkihiroSuda/multidocker)
[![Go Report Card](https://goreportcard.com/badge/github.com/AkihiroSuda/multidocker)](https://goreportcard.com/report/github.com/AkihiroSuda/multidocker)

Multidocker enables multiple Docker daemons on a single machine. No VM is required.

I made this tool so that I can easily and quickly try the latest Docker bundle.

__WARNING__: Running multiple Docker daemons is not fully supported. See also: [docker/docker#22763](https://github.com/docker/docker/pull/22763)

## How to install

    $ go install github.com/AkihiroSuda/multidocker/cmd/multidocker

Binary release will be available soon.

## Usage
Supported Docker Engine version: 1.12 or later

    $ cd ~/WORK
    $ git clone https://github.com/docker/docker.git
    $ cd docker
    $ make

    $ sudo $GOPATH/bin/multidocker create --engine-install-bundle=$HOME/WORK/docker/bundles/latest --engine-storage-driver aufs foobar
    $ eval $(multidocker env --set-client-path foobar)
    $ docker info

Note (incompatible to Docker Machine):

 * `multidocker create` copies the binaries under the bundle directory specified in `--engine-install-bundle` to `/var/lib/multidocker/hosts/foobar/bin`.
 * `eval $(multidocker env --set-client-path foobar)` sets `PATH` to `/var/lib/multidocker/hosts/foobar/bin:$PATH`.

## TODO

 * Support Docker <= 1.11
 * Support dynbinary-daemon
 * Support more installation options (e.g. `--engine-install-release=1.xx`?)
 * Support Swarm?
 * Support Windows?
 * Support UML kernel for ease of testing multiple kernels? (Maybe it should be `docker-machine-driver-uml`?)
 * Reimplement as a Docker Machine driver? (requires [docker/machine#2822](https://github.com/docker/machine/issues/2822) to be implemented)

## Related tools

 * [docker-machine-driver-dind](https://github.com/nathanleclaire/docker-machine-driver-dind)
 * [sekexe](https://github.com/jpetazzo/sekexe)

