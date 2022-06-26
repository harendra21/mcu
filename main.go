package main

import (
	"crypto/md5"
	"fmt"
	"time"
)

func main() {
	marvel_public_key := "302a803057c0480f5186b710a5454fc5"
	marvel_private_key := "ae2a810470116365dfb19025b37a5e39df4c7a1e"
	tNow := time.Now()
	tUnix := tNow.Unix()
	hash_str := fmt.Sprintf("%d%s%s", tUnix, marvel_private_key, marvel_public_key)
	data := []byte(hash_str)
	hash := md5.Sum(data)

	fmt.Println(tUnix)
	fmt.Println()
	fmt.Printf("%x", hash)
	fmt.Println()
}
