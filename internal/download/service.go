package download

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const TARGET_DIRECTORY = "data"

func SaveFiles() {
	libRegEx, e := regexp.Compile("^(data\\.json)$")
	if e != nil {
		log.Fatal(e)
	}

	filepath.Walk(".", func(filepath string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {

			jsonFile, err := os.Open(filepath)
			if err != nil {
				fmt.Println(err)
				return err
			}

			defer func(json *os.File) {
				_ = json.Close()
			}(jsonFile)

			byteValue, _ := ioutil.ReadAll(jsonFile)
			var downloads []Download

			json.Unmarshal(byteValue, &downloads)

			dataDir := strings.TrimRight(filepath, info.Name()) + "/" + TARGET_DIRECTORY
			_, err = os.Stat(dataDir)
			if err == nil {
				downloads = removeDuplicates(downloads, dataDir)
			} else {
				createDir(dataDir)
			}

			for _, download := range downloads {
				err := saveFile(download, dataDir)
				if err != nil {
					return fmt.Errorf("error saving download %s", err)
				}
			}

		}
		return nil
	})
}

func removeDuplicates(downloads []Download, dataDir string) []Download {
	dir, err := ioutil.ReadDir(dataDir)
	if err != nil {
		return nil
	}

	sort.Slice(downloads, func(i, j int) bool {
		return downloads[i].Name < downloads[j].Name
	})

	sort.Slice(dir, func(i, j int) bool {
		return dir[i].Name() < dir[j].Name()
	})

	for i, download := range downloads {
		for j, file := range dir {
			if download.Name == file.Name() {
				downloads = append(downloads[:i], downloads[i+1:]...)
				dir = append(dir[:j], dir[j+1:]...)
			}
		}
	}

	return downloads
}

func createDir(dataDir string) error {
	return os.Mkdir(dataDir, 0755)
}

func saveFile(download Download, dataDir string) error {
	response, err := http.Get(download.Url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(dataDir + "/" + download.Name)
	if err != nil {
		return err
	}

	return saveDataToFile(file, response.Body)
}

func createTempFile(fileName string) (*os.File, error) {
	return os.Create(TARGET_DIRECTORY + "/" + fileName)
}

func saveDataToFile(file io.Writer, data io.Reader) error {
	_, err := io.Copy(file, data)
	return err
}
