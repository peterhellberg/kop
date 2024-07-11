package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/peterhellberg/kop/rpc"
)

const defaultEndpoint = "http://localhost:12432/rpc/"

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	for i := range args {
		args[i] = strings.ToUpper(args[i])
	}

	list := rpc.NewListClient(rpc.New(endpoint()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(args) == 0 {
		res, err := list.Items(ctx, rpc.ItemsRequest{})
		if err != nil {
			return err
		}

		return bullets(stdout, res.Items)
	}

	switch args[0] {
	case "INGET", "NOTHING":
		_, err := list.Clear(ctx, rpc.ClearRequest{})

		return err
	case "INTE", "NO":
		res, err := list.Remove(ctx, rpc.RemoveRequest{Items: args[1:]})
		if err != nil {
			return err
		}

		return bullets(stdout, res.Items)
	default:
		res, err := list.Add(ctx, rpc.AddRequest{Items: args})
		if err != nil {
			return err
		}

		return bullets(stdout, res.Items)
	}
}

func endpoint() string {
	if endpoint := os.Getenv("KOP_ENDPOINT"); endpoint != "" {
		return endpoint
	}

	return defaultEndpoint
}

func bullets(w io.Writer, items []string) error {
	for _, item := range items {
		fmt.Fprintln(w, " -", item)
	}

	return nil
}
