#!/bin/bash

install_dependencies() {
	go get "github.com/gorilla/sessions"
	go get "golang.org/x/crypto/bcrypt"
	go get "github.com/mssola/user_agent"
	go get "github.com/gorilla/mux"
	go get "github.com/Pallinder/go-randomdata"
	go get "github.com/gorilla/handlers"
	go get "github.com/go-sql-driver/mysql"
	go get -u "github.com/mattes/migrate"
	go get "github.com/joho/godotenv"
	go get "github.com/robfig/cron"
}

bin() {
	go build
}

database_migrate() {
	export $(cat .env | xargs)
	$GOPATH/bin/migrate -url mysql://$ANA_DATABASE_USER:$ANA_DATABASE_PASSWORD@$ANA_DATABSE_HOST/$ANA_DATABASE_NAME -path ./db/migrations $1 $2 $3
}

# call first argument
$1 $2 $3 $4
