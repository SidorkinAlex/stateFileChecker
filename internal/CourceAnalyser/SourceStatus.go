package CourceAnalyser

import (
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"github.com/SidorkinAlex/stateFileChecker/internal/CliApgParser"
	"github.com/SidorkinAlex/stateFileChecker/internal/CliTextColor"
	"github.com/SidorkinAlex/stateFileChecker/internal/Encoder"
	"github.com/SidorkinAlex/stateFileChecker/internal/ManifestReader"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const fileNameIgnore = ".consistencyIgnore"
const directoryFilesChecker = ".consistency"
const analyseFileSource = ".result.lock"
const manifestFile = "manifest.json"

func CheckHashes(parser CliApgParser.CliParser) (bool, string) {
	// Read ignore file
	root := parser.Sources

	arrIgnore := createIgnoreDirList(root)

	// Open hash file
	hashFilePath := filepath.Join(root, directoryFilesChecker, analyseFileSource)
	hashFile, err := os.Open(hashFilePath)
	if err != nil {
		return false, CliTextColor.SetRedColor("Ошибка при открытии файла: " + err.Error() + "\n")
	}
	defer hashFile.Close()

	reader := csv.NewReader(hashFile)
	records, err := reader.ReadAll()
	if err != nil {
		return false, CliTextColor.SetRedColor("Ошибка при чтении файла: " + err.Error() + "\n")
	}

	manifest := ManifestReader.ManifestRead(filepath.Join(root, manifestFile))
	if manifest.Version == "" {
		return false, "manifest has empty version"
	}

	// Create hash map from file
	hashMap := make(map[string]string, len(records))
	for _, record := range records {
		if len(record) == 2 {
			decodedKey := Encoder.DecodeFromKey(record[0], manifest.Version)
			hashMap[decodedKey] = Encoder.DecodeFromKey(record[1], manifest.Version)
		}
	}

	newHashMap := make(map[string]string)
	var fileDiff HashMapDiff

	// Walk through files and check hashes
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip ignored files
		for _, value := range arrIgnore {
			if strings.Contains(path, root+"/"+value) {
				return nil
			}
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			h := sha1.New()
			if _, err := io.Copy(h, file); err != nil {
				return err
			}
			hash := fmt.Sprintf("%x", h.Sum(nil))

			// Get relative path
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			newHashMap[relPath] = hash
		}
		return nil
	})

	if err != nil {
		return false, CliTextColor.SetRedColor("Ошибка при обходе файлов: " + err.Error() + "\n")
	}

	fileDiff = FileMapDiff(hashMap, newHashMap)
	return !fileDiff.HasDiff(), fileDiff.LogString()
}

func FileMapDiff(oldHashMap, newHashMap map[string]string) HashMapDiff {
	missingFiles := []string{}
	newFiles := []string{}
	changedFiles := []string{}

	for key, oldHash := range oldHashMap {
		if newHash, exists := newHashMap[key]; !exists {
			missingFiles = append(missingFiles, key)
		} else if oldHash != newHash {
			changedFiles = append(changedFiles, key)
		}
	}

	for key := range newHashMap {
		if _, exists := oldHashMap[key]; !exists {
			newFiles = append(newFiles, key)
		}
	}

	return NewHashMapDiff(changedFiles, newFiles, missingFiles)
}

func filterEmptyStrings(arr []string) []string {
	filtered := make([]string, 0, len(arr))
	for _, str := range arr {
		if str != "" {
			filtered = append(filtered, str)
		}
	}
	return filtered
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createIgnoreDirList(rootPath string) []string {
	var arrIgnore []string

	ignoreFile := rootPath + "/" + fileNameIgnore
	ignoreData, err := os.ReadFile(ignoreFile)
	if err != nil {
		fmt.Println(err)
	}
	arrIgnore = filterEmptyStrings(strings.Split(string(ignoreData), "\n"))
	arrIgnore = append(arrIgnore, directoryFilesChecker)
	return filterEmptyStrings(arrIgnore)
}
