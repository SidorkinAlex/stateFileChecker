package main

import (
	"fmt"
	"github.com/SidorkinAlex/stateFileChecker/internal/CliApgParser"
	"github.com/SidorkinAlex/stateFileChecker/internal/CliTextColor"
	"github.com/SidorkinAlex/stateFileChecker/internal/CourceAnalyser"
	"log"
)

func main() {
	Args := CliApgParser.GetArgs()
	if Args.Action == "init" {
		CourceAnalyser.CheckHashes(Args)
	}
	log.Println(CliTextColor.SetGreenColor("success checking state app in dir " + Args.Sources))
	fmt.Println("\n")
	fmt.Println("\n")
}
