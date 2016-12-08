#!/bin/bash

bin() {
	go run ana.go
}

migrate() {
	export $(cat .env | xargs)
	$GOPATH/bin/migrate -url mysql://$ANA_DATABASE_USER:$ANA_DATABASE_PASSWORD@$ANA_DATABSE_HOST/$ANA_DATABASE_NAME -path ./db/migrations $1 $2 $3
}

# call first argument
$1 $2 $3 $4
