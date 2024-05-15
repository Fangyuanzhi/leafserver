module lf

go 1.22.1

require github.com/name5566/leaf v0.0.0-20221021105039-af71eb082cda

require github.com/gomodule/redigo v1.9.2 // indirect

require (
	github.com/gorilla/websocket v1.5.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	server v0.0.0
)

replace server => ../src/server
