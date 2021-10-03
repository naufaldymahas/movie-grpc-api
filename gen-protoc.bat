protoc --proto_path=proto --go_out=plugins=grpc:pb movie.proto
protoc --proto_path=proto --grpc-gateway_out=logtostderr=true:pb movie.proto
protoc --proto_path=proto --openapiv2_out=logtostderr=true:docs movie.proto