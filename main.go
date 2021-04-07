package main

import("fmt"
	"net/http"
)
func main() {
	fmt.Println("Hello")
	localhost := "http://localhost.com"
	_,err := http.Get(localhost)

	if err != nil{
		fmt.Println("Error")
	}

}
