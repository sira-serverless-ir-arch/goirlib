package disk

import (
	"log"
	"os"
)

func SaveFileOnDisk(folder, fileName string, buf []byte) error {
	err := os.WriteFile(folder+"/"+fileName, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}

func ReadFileOnDisk(folder, fileName string) ([]byte, error) {
	buf, err := os.ReadFile(folder + "/" + fileName)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func CreteDirIfNotExist(dirName string) {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
	}
}

func ListDirectories(path string) []string {
	var directories []string

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories
}
