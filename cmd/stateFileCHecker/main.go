package main

import (
	"github.com/SidorkinAlex/stateFileChecker/internal/CliApgParser"
	"github.com/SidorkinAlex/stateFileChecker/internal/CliTextColor"
	"github.com/SidorkinAlex/stateFileChecker/internal/CourceAnalyser"
	"log"
	"os/exec"
)

func main() {
	Args := CliApgParser.GetArgs()
	successCHeck, fileDiff := CourceAnalyser.CheckHashes(Args)
	if successCHeck {
		if len(Args.SuccessCommand) > 0 {
			exec.Command("/bin/sh", "-c", Args.SuccessCommand).Output()
		}
		log.Println(CliTextColor.SetGreenColor("the consistency of the directory has been successfully checked " + Args.Sources))
	}
	if len(Args.FailedCommand) > 0 {
		exec.Command("/bin/sh", "-c", Args.FailedCommand).Output()
	}
	log.Fatal("\n" + fileDiff)

}
