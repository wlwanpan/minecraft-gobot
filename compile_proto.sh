#!/bin/bash
protoc -I messages/ -I${GOPATH}/src --go_out=plugins=grpc:messages/ messages/cmd.proto
