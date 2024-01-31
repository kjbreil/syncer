# syncer

protoc --go_out=. \
--go-grpc_out=.  \
control/proto/*.proto

protoc -I=control/proto --go_out=. control/proto/*.proto