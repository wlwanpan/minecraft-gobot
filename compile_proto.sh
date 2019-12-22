#!/bin/bash
protoc -I services/ -I${GOPATH}/src --go_out=plugins=grpc:services/ services/*.proto
