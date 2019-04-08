# [POC]Docker Linux multicast within docker0 network

## Purpose

This repo demonstrate that it is possible to multicast among Linux docker containers.

## Prerequisite

- Docker installed

## Setup

```shell
make build
```

## Demo

Open one terminal

```shell
make run
```

Open another terminal

```shell
make run
```

You can run as many docker instance as you want. The more broadcasted messages you will received in each terminal.
