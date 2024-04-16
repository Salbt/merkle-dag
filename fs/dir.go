package fs

import (
	"log"
	merkledag "merkle-dag"
	"os"
	"path/filepath"
)

type IPFSDir struct {
	size     	uint64
	name     	string
	nodeType 	int
	it       	merkledag.DirIterator
}

type IPFSDirIterator struct {
	parentDir 	string
	nodes     	[]string
	index     	int
}

func (dir IPFSDir) Name() string {
	return dir.name
}

func (dir IPFSDir) Size() uint64 {
	return dir.size
}

func (dir IPFSDir) Type() int {
	return dir.nodeType
}

func (dir IPFSDir) It() merkledag.DirIterator {
	return dir.it
}

func (it *IPFSDirIterator) Next() bool {
	it.index++
	return it.index <= len(it.nodes)-1
}

func (it *IPFSDirIterator) Node() merkledag.Node {
	var node merkledag.Node

	path := filepath.Join(it.parentDir, it.nodes[it.index])
	
	file, err := os.Stat(path)
	if err != nil {
		log.Fatal("get file info failed: ", err)
	}
	
	if !file.IsDir() {
		node = NewFile(path)
	} else {
		node = NewDir(path)
	}

	return node
}

func NewDir(path string) *IPFSDir {
	var ipfs_dir IPFSDir
	dirIt := IPFSDirIterator{
		nodes: make([]string, 0),
		index: 0,
	}
	
	// get dir info
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal("get file Info failed", err)
	}
	
	// get all file in current dir 
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal("read file failed", err)
	}
	
	// range fileName to it.nodes
	for _, file := range files {
		dirIt.nodes = append(dirIt.nodes, file.Name())
	}
	
	// init dir
	ipfs_dir.name = fileInfo.Name()
	ipfs_dir.size = uint64(fileInfo.Size())
	ipfs_dir.nodeType = merkledag.DIR
	
	// init dirIt
	dirIt.parentDir = filepath.Join(filepath.Dir(path), ipfs_dir.name) 
	ipfs_dir.it = &dirIt			
	return &ipfs_dir
}
