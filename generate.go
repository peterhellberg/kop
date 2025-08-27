package kop

//go:generate mkdir -p ./rpc
//go:generate go run -modfile=./oto/go.mod ./oto/main.go -template template.plush -out ./rpc/gen.go -pkg rpc ./definitions
//go:generate goimports -w ./rpc/gen.go ./rpc/gen.go
