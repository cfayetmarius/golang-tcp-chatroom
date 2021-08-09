package main

import(
	"fmt"
	b64"encoding/base64"
)

func b64enc(msg string) string {
	msgenc := b64.StdEncoding.EncodeToString([]byte(msg))
	return msgenc
}	

func b64dec(msg string) string {
	msgdec := b64.StdEncoding.EncodeToString([]byte(msg))
	return msgdec
}

func main() {
	fmt.Println(b64enc("lala\n"))
	fmt.Println(b64enc("\n"))
}