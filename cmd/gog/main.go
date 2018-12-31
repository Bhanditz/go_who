package main

import (
	"github.com/mchirico/go_who/pkg"
)

func main() {

	a := pkg.App{}
	a.Initilize()
	a.Run("8080", 15, 15)

}
