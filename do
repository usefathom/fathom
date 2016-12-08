#!/bin/bash

bin() {
	go install
	sudo service ana restart
}

migration() {
	env $(cat .env | xargs) | $GOPATH/bin/migrate -url mysql://$ANA_DATABASE_USER:$ANA_DATABASE_PASSWORD@$ANA_DATABSE_HOST/$ANA_DATABASE_NAME -path ./db/migrations create $1
}

# call first argument
$1 $2 $3 $4
