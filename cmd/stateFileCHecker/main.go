package main

import (
	"github.com/SidorkinAlex/stateFileChecker/internal/CliApgParser"
	"github.com/SidorkinAlex/stateFileChecker/internal/CliTextColor"
	"github.com/SidorkinAlex/stateFileChecker/internal/CourceAnalyser"
	"log"
)

func main() {
	Args := CliApgParser.GetArgs()
	successCHeck, fileDiff := CourceAnalyser.CheckHashes(Args)
	if successCHeck {
		log.Println(CliTextColor.SetGreenColor("the consistency of the directory has been successfully checked " + Args.Sources))
	}
	log.Fatal(fileDiff)
}
