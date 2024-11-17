.PHONY: proto_gen
proto_gen:
	protoc  --proto_path=./proto --go_out=common/rpc  --go-grpc_out=common/rpc --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative   proto/*.proto


.PHONY: user_test
user_test:
	cd ./services/user && go test ./...

.PHONY: order_test
order_test:
	cd ./services/order && go test ./...



.PHONY: e2e_test
e2e_test:
	cd ./e2e && USER_SERVICE_ADDR="127.0.0.1:5002" ORDER_SERVICE_ADDR="127.0.0.1:5001" go test -v -count=1 ./...