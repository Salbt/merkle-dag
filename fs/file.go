package fs

import (
	"fmt"
	"io"
	"log"
	merkledag "merkle-dag"
	"os"
)


type IPFSFile struct {
	size     	uint64
	name     	string
	nodeType 	int
	data 		[]byte
}

func (file IPFSFile) Size() uint64 {
	return file.size
}

func (file IPFSFile) Name() string {
	return file.name
}

func (file IPFSFile) Type() int {
	return file.nodeType
}

func NewFile(path string) *IPFSFile {
	var ipfs_file IPFSFile
	
	content, err := os.Open(path)
	if err != nil {
		log.Fatal("open file falied")
	}
	defer content.Close()

	fileInfo, err := content.Stat()
	if err != nil {
		log.Fatal("get fileInfo failed", err)
	}
	
	ipfs_file.name = fileInfo.Name()
	ipfs_file.size = uint64(fileInfo.Size())
	ipfs_file.nodeType = merkledag.FILE
	
	ipfs_file.data, err = GetFileData(content)
	if err != nil {
		log.Fatal("get file data failed", err)
	}

	return &ipfs_file
}

func (file IPFSFile) Bytes() []byte {
	return file.data
}

func GetFileData(file *os.File) (data []byte, err error) {
	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)

		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("read file failed: %s", err)
		}

		if n == 0 {
			break	
		}

		data = append(data, buffer[:n]...)
	}
	return data, nil
}