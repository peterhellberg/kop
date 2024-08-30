package kop

//go:generate mkdir -p ./rpc
//go:generate go run github.com/pacedotdev/oto@v0.20.0 -template template.plush -out ./rpc/gen.go -pkg rpc ./definitions
//go:generate goimports -w ./rpc/gen.go ./rpc/gen.go
