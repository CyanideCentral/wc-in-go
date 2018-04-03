package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	if showLines {
		out += fmt.Sprintf(" | %v", ca.lines)
	}
	if showWords {
		out += fmt.Sprintf(" | %v", ca.words)
	}
	if showBytes {
		out += fmt.Sprintf(" | %v", ca.bytes)
	}
	if showChars {
		out += fmt.Sprintf(" | %v", ca.chars)
	}
	if showMaxlen {
		out += fmt.Sprintf(" | %v", ca.maxlen)
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

func isLetter(ch rune) bool {
	if 'a' <= ch && ch <= 'z' {
		return true
	}
	if 'A' <= ch && ch <= 'Z' {
		return true
	}
	return false
}

func sizeNonEmpty(slist []string) int {
	num := 0
	for _, s := range slist {
		if len(s) > 0 {
			num++
		}
	}
	return num
}

func countFile(fullpath string) *count {
	cnt := new(count)
	b, err := ioutil.ReadFile(fullpath)
	if err != nil {
		log.Println(err)
	}
	cnt.bytes = len(b)
	str := fmt.Sprintf("%s", b)
	cnt.chars = len(str)

	lineArray := strings.Split(str, "\n")
	cnt.lines = strings.Count(str, "\n")
	//TODO: tab width = 8- (pos % 8)
	maxl := 0
	for _, lstr := range lineArray {
		if len(lstr) > maxl {
			maxl = len(lstr)
		}
	}

	cnt.maxlen = maxl
	//TODO: count by appearance of consecutive letters
	words := strings.Split(str, " ")
	cnt.words = sizeNonEmpty(words)

	return cnt
}

func main() {
	args := os.Args[1:]
	i := 0
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	//Parse option flags
	for ; i < len(args); i++ {
		if args[i][0] == '-' {
			for j := 1; j < len(args[i]); j++ {
				parseParam(args[i][j : j+1])
			}
		} else {
			break
		}
	}
	//Default options
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
