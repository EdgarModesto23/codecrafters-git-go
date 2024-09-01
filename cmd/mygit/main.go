package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

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
      encoded = data.Bytes()[idx + 1:]
    }
  }
  res := fmt.Sprintf("%s", encoded)
	return res, nil
}

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		if os.Args[2] == "-p" {
			dir := os.Args[3][:2]
			file := os.Args[3][2:]
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

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
