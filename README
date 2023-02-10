# Instructions

## Building Image

```bash
docker build . --tag richie-project1
```

## Creating Bridge Network

```bash
docker network create richie-project1-network
```

## Starting Server

```bash
docker run --rm -it --network richie-project1-network --name server richie-project1 -c "servercli -address 0.0.0.0"
```

## Starting Client

```bash
docker run --rm -it --network richie-project1-network --name client richie-project1 -c "clientcli"
```

The client is interactive using the standard commands. You can learn about the available commands by typing `-h` or learn how to use a particular command by typing `<command> -h`.

To connect to the server, type:
```
c server:3000
```
Where `server` is the Docker DNS name corresponding to the name of the server container, and `3000` is the default server port.
