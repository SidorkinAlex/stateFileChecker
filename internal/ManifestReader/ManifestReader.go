package ManifestReader

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Manifest struct {
	Version string `json:"version"`
}

func ManifestRead(filePath string) Manifest {

	// Чтение содержимого файла
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Ошибка чтения файла:", err)
	}

	// Декодирование JSON данных в структуру
	var data Manifest
	err = json.Unmarshal(fileData, &data)
	if err != nil {
		log.Fatal("Ошибка декодирования JSON:", err)
	}

	return data
}
