package CliApgParser

import (
	"flag"
)

type CliParser struct {
	Sources        string
	TargetDir      string
	SuccessCommand string
	FailedCommand  string
}

//type Source struct {
//	Path string
//	Hash string
//}

func GetArgs() CliParser {

	CliParserCar := CliParser{}
	flag.StringVar(&CliParserCar.Sources, "s", "", "Sources parameter")
	flag.StringVar(&CliParserCar.SuccessCommand, "success--run", "", "command to exec if success checking")
	flag.StringVar(&CliParserCar.FailedCommand, "failed--run", "", "command to exec if failed checking")
	flag.Parse()

	return CliParserCar
}
