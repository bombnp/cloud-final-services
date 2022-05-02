#!/usr/bin/bash

connection="postgresql://postgres:postgres@localhost:5432/message_sync?sslmode=disable"
migrate -path migrations -database "$connection" up
