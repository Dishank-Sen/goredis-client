package client

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestClient(t *testing.T){
	cfg := Config{
		addr: ":5000",
	}
	client := NewClient(cfg)
	defer client.Close()

	if err := client.Set(context.Background(), "name", "Dishank"); err != nil{
		panic(err)
	}

	fmt.Println("key set")

	curr := time.Now()
	val, err := client.Get(context.Background(), "name")
	if err != nil{
		panic(err)
	}
	fmt.Println("duration: ", time.Since(curr))
	fmt.Printf("value: %s", val)
}