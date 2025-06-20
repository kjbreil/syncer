
## Control
Protobuf definition for data to be sent and GRPC definition for a server/client. Extractor creates `control.Entries`
from a struct and Injector uses that data to change a struct to match the `control.Entries` sent. GRPC definition
contains services for Push, Pull, Push/Pull and sending control messages. Everything is client based meaning the client
sets up the services, the server only acts as an endpoint responding to the clients requests. This is where the
Push/Pull service is useful because with it the client opens up a service to which the server can send data when changes
are detected rather than the client needing to blinding requesting changes.



protoc -I=control/proto --go_out=. --js_out=import_style=commonjs,binary:control/web --grpc-web_out=import_style=commonjs,mode=grpcwebtext:control/web control/proto/*.proto


protoc --js_out=import_style=commonjs,binary:. --grpc-web_out=import_style=commonjs,mode=grpcwebtext:. catalog.proto

