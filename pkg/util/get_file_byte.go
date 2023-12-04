package util

import "os"

func GetFileBytes(file *os.File) []byte {
	fileBytes := make([]byte, 0)
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileBytes = make([]byte, fileSize)

	// Read the file into fileBytes
	buffer := make([]byte, 512)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		fileBytes = append(fileBytes, buffer[:n]...)
	}

	return fileBytes
}
