package main

import "context"
import "github.com/jl3953/gotestdrabbles"

func main() {
	ctx := context.Background()
	_ = gotestdrabbles.Read(ctx)
}
