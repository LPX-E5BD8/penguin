package main

import (
	"fmt"

	"github.com/liipx/penguin/model"
)

func main() {
	info, _ := model.NewReleaseInfo(model.Version55)
	fmt.Println(info)
}
