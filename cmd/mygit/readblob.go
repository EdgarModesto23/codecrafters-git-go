package main

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type ObjectRead struct {
	command string
	flag    string
	args    []string
}

type BlobRead struct {
	flag string
	args []string
}

type TreeRead struct {
	flag string
	args []string
}

type ObjectReader interface {
	ReadBlob(f *os.File) (string, error)
}

func (b *BlobRead) ReadBlob(f *os.File) (string, error) {
	if b.flag == "-p" {
		r, err := zlib.NewReader(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var encoded []byte
		defer r.Close()
		data := new(bytes.Buffer)
		io.Copy(data, r)
		for idx, val := range data.Bytes() {
			if val == 0 {
				encoded = data.Bytes()[idx+1:]
			}
		}
		res := fmt.Sprintf("%s", encoded)
		return res, nil
	}
	return "", errors.New("Flag not recognized")
}

func (t *TreeRead) ReadBlob(f *os.File) (string, error) {
	r, err := zlib.NewReader(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data := new(bytes.Buffer)
	io.Copy(data, r)
	r.Close()
	tree, err := GetTree(data)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var res string
	if t.flag == "--name-only" {
		for _, v := range tree.objects {
			res += fmt.Sprintln(v.name)
		}
	}
	return res, nil
}

func (c *ObjectRead) GetObjectRead(cmd string) ObjectReader {
	if cmd == "cat-file" {
		return &BlobRead{flag: c.flag, args: c.args}
	}
	return &TreeRead{flag: c.flag, args: c.args}
}

func (c *ObjectRead) Exec() {
	r := c.GetObjectRead(c.command)
	var sha_position int
	if len(c.args) > 1 {
		sha_position = 1
	}
	dir := c.args[sha_position][:2]
	file := c.args[sha_position][2:]
	path := fmt.Sprintf(".git/objects/%v/%v", dir, file)
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	data, err := r.ReadBlob(f)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(data)
}
