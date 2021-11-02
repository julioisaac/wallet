#! /bin/bash

curl --location --request POST 'http://localhost:8000/create' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user" : "ben"
}'

curl --location --request POST 'http://localhost:8000/deposit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user" : "ben",
    "amount" : {
        "id": "bitcoin",
        "currency": "btc",
        "value": 0.4
    }
}'

curl --location --request POST 'http://localhost:8000/withdraw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user" : "ben",
    "amount" : {
        "id": "bitcoin",
        "currency": "btc",
        "value": 0.1
    }
}'

curl --location --request POST 'http://localhost:8000/currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id" : "usd",
    "name" : "dollar"
}'

curl --location --request POST 'http://localhost:8000/currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id" : "eur",
    "name" : "euro"
}'

curl --location --request POST 'http://localhost:8000/crypto-currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "symbol" : "btc",
    "id" : "bitcoin"
}'

curl --location --request POST 'http://localhost:8000/crypto-currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "symbol" : "eth",
    "id" : "ethereum"
}'

curl --location --request POST 'http://localhost:8000/crypto-currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "symbol" : "ada",
    "id" : "cardano"
}'

curl --location --request POST 'http://localhost:8000/crypto-currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "symbol" : "doge",
    "id" : "dogedoin"
}'

curl --location --request POST 'http://localhost:8000/crypto-currency' \
--header 'Content-Type: application/json' \
--data-raw '{
    "symbol" : "xrp",
    "id" : "xrp"
}'