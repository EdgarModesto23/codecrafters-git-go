package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"os"
)

type BlobWrite struct {
	flag     string
	filepath string
}

func SplitObjName(sha string) (dir string, fname string) {
	d := sha[:2]
	fn := sha[2:]

	return d, fn
}

func CreateFiles(dir string, fname string) (string, error) {
	fpath := fmt.Sprintf(".git/objects/%v", dir)

	err := os.Mkdir(fpath, 0755)
	if err != nil {
		return "", err
	}
	f, err := os.OpenFile(fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fullpath := fpath + fmt.Sprintf("/%v", fname)

	return fullpath, nil
}

func WriteGitObjectFile(fpath string, data []byte) error {
	var buf bytes.Buffer
	zwritter := zlib.NewWriter(&buf)
	_, err := zwritter.Write(data)
	zwritter.Close()
	if err != nil {
		return err
	}
	err = os.WriteFile(fpath, buf.Bytes(), 0755)
	if err != nil {
		return err
	}
	return nil
}

func (w *BlobWrite) GetFileContent() ([]byte, error) {
	f_content, err := os.ReadFile(w.filepath)
	if err != nil {
		return nil, nil
	}
	return f_content, nil
}

func BytestoSHA(content []byte) (string, string) {
	header := fmt.Sprintf("blob %v\x00", len(content))
	fsha := append([]byte(header), content...)
	h := sha1.New()
	h.Write(fsha)
	bs := fmt.Sprintf("%x", h.Sum(nil))
	return bs, header
}

func (w *BlobWrite) Exec() {
	if w.flag == "-w" {
		f_content, err := w.GetFileContent()
		if err != nil {
			panic(err)
		}
		bs, header := BytestoSHA(f_content)
		d, fn := SplitObjName(bs)
		fname, err := CreateFiles(d, fn)
		if err != nil {
			panic(err)
		}
		WriteGitObjectFile(fname, append([]byte(header), f_content...))
		fmt.Println(bs)
	}
}
