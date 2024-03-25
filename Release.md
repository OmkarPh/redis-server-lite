
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
| SET     	| SET <key> <value> 	| redis-cli SET name Mark      	|   	|
| GET     	| GET <key>         	| redis-cli GET name           	|   	|
| INCR    	| INCR key          	| redis-cli INCR age           	|   	|
| DECR    	| DECR key          	| redis-cli DECR age           	|   	|
| EXISTS    	| EXISTS key [key ...]           	| redis-cli EXISTS name age           	|   	|
| EXPIRE    	| EXPIRE key seconds         	| redis-cli EXPIRE name 20           	|   	|
| TTL    	| TTL key          	| redis-cli TTL key           	|   	|
| DEL    	| DEL key [key ...]           	| redis-cli DEL name age           	|   	|
| TYPE    	| TYPE key          	| redis-cli TYPE name           	|   	|
| PING    	| PING              	| redis-cli PING               	|   	|
| ECHO    	| ECHO <message>    	| redis-cli ECHO "Hello world" 	|   	|
