package merkledag_test

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	merkledag "merkle-dag"
	"merkle-dag/fs"
	kvstore "merkle-dag/kvStore"
	"os"
	"testing"
)

// dir1 目录通过merkledag计算得到的hash值
const ROOTHASH = "d8cc3b28f3ecc2906e8265c759ecd8181baa8bc347f929081d45a3158983c0f5"

var store *kvstore.LevelDB

func TestMain(m *testing.M) {
	CreateTestFile()                 // 创建一个文件夹
	store = kvstore.NewLevelDB("db") // 创建leveldb

	exitCode := m.Run()

	RemoveDir("dir1") // 删除创建的文件夹
	store.Close()
	RemoveDir("db") // 删除leveldb

	os.Exit(exitCode)
}

func TestNewNode(t *testing.T) {

	intput := "dir1"

	var node merkledag.Node = *fs.NewNode(intput)

	if node.Type() != merkledag.DIR || node.Name() != "dir1" {
		t.Error(`NewNode(dir1) :file to node failed`)
	}

}

func TestDagAdd(t *testing.T) {

	dir1 := *fs.NewNode("dir1")
	hasher := sha256.New()

	key := merkledag.Add(store, dir1, hasher)
	keyHex := hex.EncodeToString(key)

	if key == nil {
		t.Error(`merkledag.Add(db, dir1, hasher): add object to db failed`)
	}

	if keyHex != ROOTHASH {
		t.Error(`merkledag.Add(db, dir1, hasher): add function get hash not match rootHash`)
	}

}

func TestDagHash2File(t *testing.T) {

	dir1 := *fs.NewNode("dir1")
	hasher := sha256.New()
	merkledag.Add(store, dir1, hasher)

	hash2Bytes, _ := hex.DecodeString(ROOTHASH)
	data := merkledag.Hash2File(store, hash2Bytes, "dir3/dir4/c.txt", nil)
	if data == nil {
		t.Error(`merkedag.Hash2File(db, file, hasher): Hash2File get data is nil`)
	}
	if string(data) != "i am a txt file" {
		t.Error(`merkedag.Hash2File(db, file, hasher): Hash2File get worng data`)
	}
}

func CreateTestFile() {
	/*
		创建的文件结构：
		-dir1
		 -dir2
		  -a.go
		  -b.js
		 -dir3
		  -dir4
		   -c.txt
	*/
	dirs := []string{"dir1/dir2", "dir1/dir3/dir4"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal("create dir failed", err)
		}
	}

	fileNames := []string{"a.go", "b.js", "c.txt"}
	content := []string{
		`package main
		import "fmt"
		func main() {
			fmt.Println("Hello World")
		}
		`,
		`function main() {
			console.log("hello world")
		}`,
		`i am a txt file`,
	}
	for index, fileName := range fileNames {
		path := dirs[0]
		if fileName == "c.txt" {
			path = dirs[1]
		}

		file, err := os.Create(path + "/" + fileName)
		if err != nil {
			log.Fatal("create file failed")
		}
		defer file.Close()

		_, err = file.WriteString(content[index])
		if err != nil {
			log.Fatal("write data to file failed")
		}
	}

}

func RemoveDir(filePath string) {

	_, err := os.Stat(filePath)
	if err != nil {
		log.Fatal("dir was not created")
	}

	err = os.RemoveAll(filePath)
	if err != nil {
		log.Fatal("remove file failed")
	}
}
