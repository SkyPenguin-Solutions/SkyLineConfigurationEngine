package main

import (
	"fmt"
	"os"

	ENGINE "github.com/SkyPenguin-Solutions/SkyLineConfigurationEngine/Engine/Backend"
)

func init() {
	ENGINE.Clear()
	ENGINE.OutputBanner()
	fmt.Print("\n\n")
	ENGINE.OutputBoxOCode("Example.slc")
}

func main() {
	ENGINE.Start(os.Stdin, os.Stdout)
}
