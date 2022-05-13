#!/usr/bin/bash

connection="postgresql://postgres:postgres@localhost:5432/cloud-final?sslmode=disable"
migrate -path migrations -database "$connection" up
