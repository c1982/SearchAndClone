package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

const FILE_NAME = "./domains.txt"

func main() {
	source := flag.String("source", "", "C:\\Backups")
	dest := flag.String("destination", "", "C:\\vhosts\\{DOMAIN}\\http")
	postfix := flag.String("postfix", "wwwroot", "")

	flag.Parse()

	SearchAndCopy(source, dest, postfix)
}

func SearchAndCopy(source *string, dest *string, postfix *string) {
	file, err := os.Open(FILE_NAME)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		name := scanner.Text()
		sourcePath := searchSourcePath(name, *source)

		if sourcePath != "" {
			sourcePath = sourcePath + "\\" + *postfix
			destPath := strings.Replace(*dest, "{DOMAIN}", name, -1)

			log.Println(sourcePath, " > ", destPath)

			if isDirectoryExists(sourcePath) && isDirectoryExists(destPath) {

				copyResult, copyError := roboCopy(sourcePath, destPath)

				if copyError != nil {
					fmt.Println("Copy Error: ", copyError)
				} else {
					if copyResult {
						fmt.Println("Success:", name)
					}
				}
			} else {
				log.Println("Directory not found: ", name)
			}

		} else {
			log.Println("Source folder not found: ", name)
		}

		fmt.Println("")
	}
}

func roboCopy(source string, destination string) (bool, error) {

	result := true
	cmd := exec.Command("robocopy.exe", source+"\\", destination+"\\", "/MIR", "/MOVE")

	//log.Println(strings.Join(cmd.Args, ";"))

	_, err := cmd.Output()

	if err != nil {
		result = false
	}

	return result, err
}

func searchSourcePath(name string, sourceFolder string) string {

	var findingFolder string

	files, err := ioutil.ReadDir(sourceFolder)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, f := range files {
		if f.IsDir() {
			findingFolder := fmt.Sprintf("%v\\%v", sourceFolder, f.Name())
			if name == f.Name() {
				return findingFolder
			} else {
				findingFolder = searchSourcePath(name, findingFolder)
				if findingFolder != "" {
					return findingFolder
				}
			}
		}
	}

	return findingFolder
}

func isDirectoryExists(path string) bool {
	finfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return finfo.IsDir()
}
