package client

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSet(t *testing.T){
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

func TestMSet(t *testing.T){
	cfg := Config{
		addr: ":5000",
	}
	client := NewClient(cfg)
	defer client.Close()

	m := make(map[string]string)

	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("key%d", i)
		v := fmt.Sprintf("value %d", i)
		m[k] = v
	}

	if err := client.MSet(context.Background(), m); err != nil{
		panic(err)
	}

	fmt.Println("keys set")

	curr := time.Now()
	val, err := client.Get(context.Background(), "key88")
	if err != nil{
		panic(err)
	}
	fmt.Println("duration: ", time.Since(curr))
	fmt.Printf("value: %s", val)
}