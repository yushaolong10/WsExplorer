#!/bin/bash
#
###################################################
# auto generate protobuf
# file description:
#	 - ws.proto  //pb server
#
####################################################


# generate pb
protoc *.proto --go_out=plugins=grpc:./