package CourceAnalyser

import "github.com/SidorkinAlex/stateFileChecker/internal/CliTextColor"

type HashMapDiff struct {
	changedFile []string
	newFile     []string
	deletedFile []string
}

func NewHashMapDiff(changedFile []string, newFile []string, deletedFile []string) HashMapDiff {
	return HashMapDiff{
		changedFile,
		newFile,
		deletedFile,
	}
}
func (h *HashMapDiff) LogString() string {
	logString := ""
	if len(h.newFile) > 0 {
		logString += CliTextColor.SetGreenColor("Add new File:\n")
		for _, filePath := range h.newFile {
			logString += CliTextColor.SetGreenColor(filePath)
		}
	}
	if len(h.changedFile) > 0 {
		logString += CliTextColor.SetYellowColor("ChangedFile File:\n")
		for _, filePath := range h.changedFile {
			logString += CliTextColor.SetYellowColor(filePath)
		}
	}
	if len(h.deletedFile) > 0 {
		logString += CliTextColor.SetRedColor("Deleted File:\n")
		for _, filePath := range h.deletedFile {
			logString += CliTextColor.SetRedColor(filePath)
		}
	}
	return logString
}
func (h *HashMapDiff) HasDiff() bool {
	if len(h.newFile) > 0 || len(h.changedFile) > 0 || len(h.deletedFile) > 0 {
		return true
	}
	return false
}
