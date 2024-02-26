# To run these commands
# 1. Install protoc compiler
#
# (on Windows)
# choco install protoc
#
# (on Ubuntu)
# apt install protobuf-compiler
#
# 2. Install plugins
#
# go install github.com/golang/protobuf/protoc-gen-go@latest
# go install github.com/asynkron/protoactor-go/protobuf/protoc-gen-go-grain@latest
#
#
# (make sure go installation path is on your PATH)

protoc --go_out=. --go_opt=paths=source_relative \
         --go-grain_out=. --go-grain_opt=paths=source_relative messages.proto
