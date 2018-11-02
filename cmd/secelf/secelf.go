package main

import (
	"os"

	"github.com/moznion/secelf"
)

func main() {
	secelf.Run(os.Args[1:])
}
