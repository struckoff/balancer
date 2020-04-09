module github.com/struckoff/kvrouter

go 1.14

require (
	github.com/go-delve/delve v1.4.0
	github.com/golang/protobuf v1.3.5
	github.com/hashicorp/consul/api v1.4.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/struckoff/SFCFramework v0.0.0-20200219142622-9ea20220c1d7
	google.golang.org/grpc v1.28.1
)

replace github.com/struckoff/SFCFramework => /home/struckoff/Projects/Go/src/github.com/struckoff/SFCFramework
