# GORedis Client

This is a minimal client library for goredis server.

## Installation

```bash
go get github.com/Dishank-Sen/goredis-client@latest
```

## Quickstart

```bash
package main

func main(){
    cfg := Config{
		addr: ":5000",
	}
	client := NewClient(cfg)
	defer client.Close()

	if err := client.Set(context.Background(), "name", "Dishank"); err != nil{
		panic(err)
	}

	fmt.Println("key set")

	val, err := client.Get(context.Background(), "name")
	if err != nil{
		panic(err)
	}
	fmt.Printf("value: %s", val)
}
```

Expected Output

```bash
key set
value: Dishank
```