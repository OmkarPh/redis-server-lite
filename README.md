
# A Lite redis-server implementation in golang

### Prerequisites

- Go 1.11 or later

#### Setup

```bash
# Clone this repository
$ git clone https://github.com/OmkarPh/redis-lite.git

# Go into the server directory
$ cd redis-lite/server/cmd

# Run the server
$ go run .
```

## Supported redis-cli commands

| Command 	| Syntax            	| Example                      	|   	|
|---------	|-------------------	|------------------------------	|---	|
| set     	| set <key> <value> 	| redis-cli set name Mark      	|   	|
| get     	| get <key>         	| redis-cli get name           	|   	|
| incr    	| incr key          	| redis-cli incr age           	|   	|
| ping    	| ping              	| redis-cli ping               	|   	|
| echo    	| echo <message>    	| redis-cli echo "Hello world" 	|   	|

## redis-benchmarks support upcoming ...

