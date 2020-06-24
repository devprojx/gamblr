package lib

import (
	"io/ioutil"
	"log"
	"os"
)

var filePath = "./settings.json"

func SaveSettings(settings string) {
	err := ioutil.WriteFile(filePath, []byte(settings), 0644)
	if err != nil {
		log.Printf("[error] unable to wite to settings: %s", err)
	}
}

func LoadSettings() string {
	// detect if file exists
	_, err := os.Stat(filePath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filePath)
		if isError(err) {
			log.Printf("[error] unable to create settings.json: %s", err)
		}
		defer file.Close()
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("[error] unable to open settings: %s", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("[error] unable to read settings: %s", err)
	}

	return string(data)
}
