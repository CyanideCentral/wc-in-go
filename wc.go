package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ryanuber/columnize"
)

type count struct {
	bytes  int
	chars  int
	lines  int
	maxlen int
	words  int
}

func (ca *count) add(cb *count) {
	ca.bytes += cb.bytes
	ca.chars += cb.chars
	ca.lines += cb.lines
	if ca.maxlen < cb.maxlen {
		ca.maxlen = cb.maxlen
	}
	ca.words += cb.words
}

func (ca *count) toString() (out string) {
	if ca == nil {
		ca = new(count)
	}
	if showBytes {
		out += fmt.Sprintf("%v", ca.bytes)
	}
	if showChars {
		out += fmt.Sprintf(" | %v", ca.chars)
	}
	if showLines {
		out += fmt.Sprintf(" | %v", ca.lines)
	}
	if showMaxlen {
		out += fmt.Sprintf(" | %v", ca.maxlen)
	}
	if showWords {
		out += fmt.Sprintf(" | %v", ca.words)
	}
	return out
}

//-c
var showBytes = false

//-m
var showChars = false

//-l
var showLines = false

//-L
var showMaxlen = false

//-w
var showWords = false

//Num of valid params
var numParams = 0

func parseParam(param string) {
	switch param {
	case "c":
		showBytes = true
	case "L":
		showMaxlen = true
	case "l":
		showLines = true
	case "m":
		showChars = true
	case "w":
		showWords = true
	default:
		log.Fatal("invalid option -- '" + param + "'")
	}
	numParams++
}

func countFile(fullpath string) *count {
	cnt := new(count)
	cnt.chars = 200
	return cnt
}

func main() {
	args := os.Args[1:]
	i := 0
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for ; i < len(args); i++ {
		if args[i][0] == '-' {
			for j := 1; j < len(args[i]); j++ {
				parseParam(args[i][j : j+1])
			}
		} else {
			break
		}
	}
	if numParams == 0 {
		showLines = true
		showWords = true
		showBytes = true
	}
	countMap := map[string]*count{}
	for ; i < len(args); i++ {
		if args[i][0] == '/' {
			countMap[args[i]] = countFile(args[i])
		} else {
			path := wd + "/" + args[i]
			countMap[args[i]] = countFile(path)
		}
	}
	total := new(count)
	output := []string{}
	for fn, cnt := range countMap {
		out := cnt.toString()
		out += " | " + fn
		output = append(output, out)
		total.add(cnt)
	}
	if len(output) > 1 {
		output = append(output, total.toString()+" | total")
	}
	result := columnize.SimpleFormat(output)
	fmt.Println(result)
}
