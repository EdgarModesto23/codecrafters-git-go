package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

type BlobRead struct {
	args []string
}

func ReadBlob(f *os.File) (string, error) {
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

func (c *BlobRead) Exec() {
	if c.args[0] == "-p" {
		dir := c.args[1][:2]
		file := c.args[1][0:]
		path := fmt.Sprintf(".git/objects/%v/%v", dir, file)
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data, err := ReadBlob(f)
		if err != nil {
			panic(err)
		}
		fmt.Print(data)
	}
}
