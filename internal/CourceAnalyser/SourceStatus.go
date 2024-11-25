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
	// Чтение файла с хешами
	root := parser.Sources
	ignoreFile := root + "/" + fileNameIgnore
	ignoreData, err := os.ReadFile(ignoreFile)
	if err != nil {
		return false, err.Error()
	}
	var arrIgnore []string
	arrIgnore = strings.Split(string(ignoreData), "\n")
	arrIgnore = filterEmptyStrings(arrIgnore)
	hashFile, err := os.Open(root + "/" + directoryFilesChecker + "/" + analyseFileSource)
	if err != nil {
		return false, CliTextColor.SetRedColor("Ошибка при открытии файла: " + err.Error() + "\n")
	}
	defer hashFile.Close()

	reader := csv.NewReader(hashFile)
	records, err := reader.ReadAll()
	if err != nil {
		return false, CliTextColor.SetRedColor("Ошибка при чтении файла: " + err.Error() + "\n")
	}
	manifest := ManifestReader.ManifestRead(root + "/" + manifestFile)
	if manifest.Version == "" {
		return false, "manifest has empty version"
	}
	// Создание карты хешей из файла
	hashMap := make(map[string]string)

	for _, record := range records {
		if len(record) == 2 {
			if !fileExists(root + "/" + Encoder.DecodeFromKey(record[0], manifest.Version)) {
				return false, CliTextColor.SetRedColor("error checking application state in dir " + root)
			}
			hashMap[Encoder.DecodeFromKey(record[0], manifest.Version)] = Encoder.DecodeFromKey(record[1], manifest.Version)
		}
	}
	newHashMap := make(map[string]string)
	var fileDiff HashMapDiff
	// Проверка хешей файлов
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Проверка игнорируемых файлов
		for _, value := range arrIgnore {
			if strings.Contains(path, value) {
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

			// Получение относительного пути от сканируемой директории
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			newHashMap[relPath] = hash
		}
		fileDiff = FileMapDiff(hashMap, newHashMap)
		return nil
	})

	return !fileDiff.HasDiff(), fileDiff.LogString()
}

func FileMapDiff(oldHashMap map[string]string, newHashMap map[string]string) HashMapDiff {

	missingFiles := []string{}
	newFiles := []string{}
	changedFiles := []string{}

	for key, oldHash := range oldHashMap {
		newHash, exists := newHashMap[key]
		if !exists {
			// The file is missing in the new hashmap
			missingFiles = append(missingFiles, key)
		} else if oldHash != newHash {
			// The file has changed
			changedFiles = append(changedFiles, key)
		}
	}

	// Checking for new files
	for key := range newHashMap {
		if _, exists := oldHashMap[key]; !exists {
			newFiles = append(newFiles, key)
		}
	}
	return NewHashMapDiff(changedFiles, newFiles, missingFiles)
}

func filterEmptyStrings(arr []string) []string {
	filtered := make([]string, 0)

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
