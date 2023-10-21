# experiments-params

- Retrieves experiment data from Configuration Data base.

![experiments-params diagram](diagrams/experiments-params.drawio.png)

## How to run

### Start sql container and run queries located in db_schema / creation_schema.sql

```
docker run --name db \
    -e MYSQL_ROOT_PASSWORD=testroot \
    -e MYSQL_USER=jackpotian \
    -e MYSQL_PASSWORD=test \
    -e MYSQL_DATABASE=experiment_db \
    -p 3307:3306 \
    -d mysql:8.0
```

#### Troubleshooting

If you get a connection error with DBeaver, please use the following connection string:
```
jdbc:mysql://localhost:3307/experiment_db?allowPublicKeyRetrieval=true&useSSL=false
```

### Set the following env variables or ask for the .env file

```
DB_HOST=localhost
DB_PORT=3307
DB_USER=jackpotian
DB_PASSWORD=test
DB_NAME=experiment_db
```

### Build and run
```
go build
go run .
```

## Explore

```
http://localhost:8091/swagger/index.html
```

## Update swagger

```
swag init
```

## Build Docker image

```
docker build -t experiment-params .
docker run -p 8091:8091 experiment-params
```

## Example experiment creation

```
curl -X 'POST' \
  'http://localhost:8091/api/v1/experiment' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "arms": [
    {
      "name": "picture_1"
    },
    {
      "name": "picture_2"
    }
  ],
  "experiment_id": "test-xp",
  "model_parameters": {
    "input_features": [
      "age"
    ],
    "model_type": "classification",
    "output_classes": [
      "0", "1"
    ]
  },
  "parameters": "{}",
  "policy_type": "epsilon_greedy"
}'
```