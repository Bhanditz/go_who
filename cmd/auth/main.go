package main

import (
	"github.com/mchirico/go_who/configure"
	"github.com/mchirico/go_who/pkg"
	"os/user"
)

func main() {

	a := pkg.App{}
	a.Initilize()
	usr, _ := user.Current()
	file := usr.HomeDir + "/.secretHarvest"
	s, _ := configure.GetSecret(file)
	a.InitSS(&s)

	a.Run("4591", 15, 15)

}
