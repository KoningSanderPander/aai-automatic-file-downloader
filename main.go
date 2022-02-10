package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {

	file, err := os.Open("day-4/data.txt")
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	spaceSplitter := regexp.MustCompile(`\s`)

	var downloads []File
	for scanner.Scan() {
		d := spaceSplitter.Split(scanner.Text(), -1)
		downloads = append(downloads, File{Url: d[0], Name: d[1]})
	}

	err = os.Mkdir("data", 0755)
	if err != nil {
		log.Fatal(err)
	}

	for _, download := range downloads {
		Download(download)
	}

	err = file.Close()
	check(err)

}

type File struct {
	Url  string
	Name string
}

func Download(download File) {
	resp, err := http.Get(download.Url)
	check(err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	out, err := os.Create("data/" + download.Name)
	check(err)

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	check(err)

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
