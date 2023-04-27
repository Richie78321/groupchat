# Distributed Groupchat

## Design

[Design Doc 1](https://docs.google.com/document/d/1WtNcXEWKXn48iFw-FbrP8NJmXR1D3_HMcLqcgoYEg2k/edit?usp=sharing) \
[Design Doc 2](https://docs.google.com/document/d/12x_JduDuIfentusbBKKCICCVh_vfmNMYNwu4fKb7gNA/edit?usp=sharing)

## Running via the Test Script

From the root directory of the project

### Build the Image

```bash
docker build . --tag cs2510_p2
```

### Run the Test Script

```bash
python ./test/test_p2.py init
```

This will start the servers.

### Start the Client

```bash
docker run --rm -it --network cs2510 --name client cs2510_p2 -c "clientcli" 
```

The client is interactive using the standard commands. You can learn about the available commands by typing `-h` or learn how to use a particular command by typing `<command> -h`.

To connect to one of the servers, type:
```
c <hostname>:3000
```

Where `<hostname>` is the name of a server container, which is `cs2510_server` followed by an integer.

## Running Manually

From the root directory of the project

### Build the image

```bash
docker build . --tag richie-project1
```

### Create the bridge network

```bash
docker network create richie-project1-network
```

### Start the Server(s)

To get information on the commandline arguments:
```bash
docker run --rm -it --network richie-project1-network richie-project1 -c "servercli -h"
```

#### Specifying Peers

If you plan to run multiple servers, you must specify the peers for each server with the `-p` argument.

For example if there are three servers:
* `localhost:3000`
* `localhost:3001`
* `localhost:3002`

Then you would start the servers by running:
```bash
<docker run> -c "servercli --address 0.0.0.0 --port 3000 -p 1:localhost:3000 -p 2:localhost:3001 -p 3:localhost:3002 1" # Starts a server on localhost:3000 with ID 1 and two peers
<docker run> -c "servercli --address 0.0.0.0 --port 3001 -p 1:localhost:3000 -p 2:localhost:3001 -p 3:localhost:3002" 2
<docker run> -c "servercli --address 0.0.0.0 --port 3002 -p 1:localhost:3000 -p 2:localhost:3001 -p 3:localhost:3002" 3
```

### Starting the Client

```bash
docker run --rm -it --network richie-project1-network --name client richie-project1 -c "clientcli"
```
