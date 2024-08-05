package CliTextColor

func SetGreenColor(text string) string {
	return "\033[32m" + text + "\033[0m\n"
}

func SetYellowColor(text string) string {
	return "\033[33m" + text + "\033[0m\n"
}

func SetRedColor(text string) string {
	return "\033[31m" + text + "\033[0m\n"
}
