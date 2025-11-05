package main

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
)

func main() {
	u := uuid.Must(uuid.NewV7())
	fmt.Println(u)
}
