package main

import (
	omi "github.com/stormi-li/omiv1"
)

func main() {
	c := omi.NewClient(&omi.Options{
		Addr:     "localhost:6379",
		CacheDir: "cache",
	})
	c.Register.RegisterAndServe("test", "fdsf:8998", func(port string) { select {} })
}
