package fs

import (
	merkledag "merkle-dag"
	"os"
)

func NewNode(path string) *merkledag.Node {
	var node merkledag.Node
	
	info, _ := os.Stat(path)
	if !info.IsDir() {
		node = NewFile(path)
	} else {
		node = NewDir(path)
	}

	return &node
}