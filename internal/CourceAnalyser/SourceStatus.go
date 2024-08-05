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
	"log"
	"os"
	"path/filepath"
	"strings"
)

const fileNameIgnore = ".consistencyIgnore"
const directoryFilesChecker = ".consistency"
const analyseFileSource = ".result.lock"
const manifestFile = "manifest.json"

func CheckHashes(parser CliApgParser.CliParser) {
	// Чтение файла с хешами
	root := parser.Sources
	ignoreFile := root + "/" + fileNameIgnore
	ignoreData, err := os.ReadFile(ignoreFile)
	if err != nil {
		fmt.Println(err)
	}
	var arrIgnore []string
	arrIgnore = strings.Split(string(ignoreData), "\n")
	arrIgnore = filterEmptyStrings(arrIgnore)
	hashFile, err := os.Open(root + "/" + directoryFilesChecker + "/" + analyseFileSource)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла: %v\n", err)
		return
	}
	defer hashFile.Close()

	reader := csv.NewReader(hashFile)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}
	manifest := ManifestReader.ManifestRead(root + "/" + manifestFile)
	if manifest.Version == "" {
		log.Fatalln("manifest has empty version")
		return
	}
	// Создание карты хешей из файла
	hashMap := make(map[string]string)
	for _, record := range records {
		if len(record) == 2 {
			if !fileExists(root + "/" + Encoder.DecodeFromKey(record[0], manifest.Version)) {
				log.Fatalln(CliTextColor.SetRedColor("error checking application state in dir " + root))
			}
			hashMap[Encoder.DecodeFromKey(record[0], manifest.Version)] = Encoder.DecodeFromKey(record[1], manifest.Version)
		}
	}
	// Проверка хешей файлов
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, value := range arrIgnore {
			if strings.Contains(path, string(value)) {
				return nil
			}
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

			// Сравнение хешей
			if storedHash, ok := hashMap[relPath]; ok {
				if storedHash != hash {
					log.Fatalln(CliTextColor.SetRedColor("file " + relPath + " is changed"))
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка при обходе папки: %v\n", err)
	}
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
