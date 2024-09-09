package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type TreeObject struct {
	mode    string
	name    string
	sha     string
	objtype string
}

type Tree struct {
	size    int
	objects []TreeObject
}

func GetTreeSize(data bytes.Buffer) (int, int, error) {
	var header string
	offset := 0
	for idx, val := range data.Bytes() {
		if val == 0 {
			header += data.String()[:idx]
			offset += idx
			break
		}
	}
	_, a, _ := strings.Cut(header, " ")
	result, err := strconv.Atoi(a)
	if err != nil {
		return -1, 0, err
	}

	return result, offset, nil
}

func GetTree(data *bytes.Buffer) (Tree, error) {
	size, offset, err := GetTreeSize(*data)
	if err != nil {
		return Tree{}, err
	}
	res := Tree{size: size, objects: []TreeObject{}}
	buf := new(bytes.Buffer)
	buf.Write(data.Bytes()[offset+1:])
	res.PushTreeObjs(*buf)
	return res, nil
}

func (t Tree) String() string {
	var res string
	for _, val := range t.objects {
		entry := fmt.Sprintf("%s %s %s   %s\n", val.mode, val.objtype, val.sha, val.name)
		res += entry
	}
	return res
}

func (t *Tree) PushTreeObjs(data bytes.Buffer) {
	dbytes := data.Bytes()
	n := false
	count := 0
	mode := []byte{}
	name := []byte{}
	SHA := []byte{}

	for i := 0; i < len(dbytes); i++ {
		if dbytes[i] == 0 { // Check for null byte
			count = 0
			for count < 20 && i+1 < len(dbytes) { // Ensure i is within bounds
				i += 1
				SHA = append(SHA, dbytes[i])
				count++
			}
			n = false
			count = 0
			t.objects = append(
				t.objects,
				TreeObject{mode: string(mode), name: string(name), sha: string(SHA)},
			)
			mode = []byte{}
			name = []byte{}
			SHA = []byte{}
			continue
		}
		// Extract mode
		for i < len(dbytes) && !n {
			mode = append(mode, dbytes[i])
			i += 1
			if dbytes[i] == 32 {
				n = true
				i += 1
				break
			}
			if i >= len(dbytes) { // Check bounds
				break
			}
		}
		// Extract name
		for i < len(dbytes) && n {
			if dbytes[i+1] == 0 {
				name = append(name, dbytes[i])
				break
			}
			name = append(name, dbytes[i])
			i += 1
			if i >= len(dbytes) { // Check bounds
				break
			}
		}
	}
}
