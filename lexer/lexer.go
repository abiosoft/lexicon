// Copyright 2010 Abiola Ibrahim. All rights reserved.
// CSC 332 Lexicon Assignment
// Analysis of Java Programming Language using Go (http://golang.org)
package lexer

import (
	"container/vector"
	"io/ioutil"
	"strutils"
	"bufio"
	"strconv"
	"os"
	"strings"
	"fmt"
)

var (
	keywords, separators, operators vector.StringVector
)

const (
	keywordsFile   = "java/keywords"
	operatorsFile  = "java/operators"
	separatorsFile = "java/separators"
)

func loadFile(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		println(err.String())
	}
	return string(data)
}

func readLine(reader *bufio.Reader) (str string, err os.Error) {
	str, err = reader.ReadString(byte('\n'))
	return
}

func init() {
	tmp := loadFile(keywordsFile)
	keywords = strings.Fields(tmp)
	tmp = loadFile(operatorsFile)
	operators = strings.Fields(tmp)
	tmp = loadFile(separatorsFile)
	separators = strings.Fields(tmp)
}

type Lexicon struct {
	keywords, separators, operators, literals, identifiers vector.StringVector
	buffer                                                 *strutils.StringBuffer
}

func (this *Lexicon) Init(filename string) { //initialize and remove comments
	this.buffer = strutils.NewStringBuffer("")
	file, err := os.Open(filename, os.O_RDONLY, 0666)
	if err != nil {
		println(err.String())
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	commentStart := [2]string{"/*", "//"}
	commentEnd := [2]string{"*/", "\n"}
	for {
		line, err := readLine(reader)
		index := strings.Index(line, commentStart[0])
		if index != -1 {
			this.buffer.Append("\n" + line[0:index])
			for {
				line, err = readLine(reader)
				index := strings.Index(line, commentEnd[0])
				if index != -1 {
					this.buffer.Append("\n" + line[0:index])
					break
				}
			}
			continue
		}
		index = strings.Index(line, commentStart[1])
		if index != -1 {
			this.buffer.Append("\n" + line[0:index])
		} else {
			this.buffer.Append("\n" + line)
		}
		if err != nil {
			break
		}
	}
}

func (this *Lexicon) FilterStrings() {
	tokenizer := strutils.NewStrTokensDelims(this.buffer.String(), "\"", true)
	buf := strutils.NewStringBuffer("")
	for tokenizer.HasMoreTokens() {
		next := tokenizer.NextToken()
		if next == "\"" {
			tmp := "\""
			for tokenizer.HasMoreTokens() {
				next = tokenizer.NextToken()
				if next == "\"" {
					if len(tmp) > 0 && string(tmp[len(tmp)-1]) == "\\" {
						tmp += next
					} else {
						tmp += next
						break
					}
				} else {
					tmp += next
				}
			}
			this.literals.Push(tmp)
		} else {
			buf.Append(next + " ")
		}
	}
	this.buffer = buf
}

func (this *Lexicon) AddSpaces() {
	tokenizer := strutils.NewStrTokens(this.buffer.String())
	buf := strutils.NewStringBuffer("")
	for tokenizer.HasMoreTokens() {
		next := strutils.NewStringBuffer(tokenizer.NextToken())
		for _, op := range operators {
			count := strings.Count(next.String(), op)
			index := 0
			for i := 0; i < count; i++ {
				index := strings.Index(next.String()[index:], op)
				next.Replace(index, index+len(op), " "+op+" ")
				index += 2
			}
		}
		for _, sp := range separators {
			count := strings.Count(next.String(), sp)
			index := 0
			for i := 0; i < count; i++ {
				index := strings.Index(next.String()[index:], sp)
				next.Replace(index, index+len(sp), " "+sp+" ")
				index += 2
			}
		}
		buf.Append(next.String() + " ")
	}
	this.buffer = buf
}

func (this *Lexicon) Analyze() {
	tokenizer := strutils.NewStrTokens(this.buffer.String())
outer:
	for tokenizer.HasMoreTokens() {
		token := tokenizer.NextToken()
		//keywords, separators, operators, literals, identifiers
		for _, k := range keywords {
			if k == token {
				this.keywords.Push(token)
				continue outer
			}
		}
		for _, s := range separators {
			if s == token {
				this.separators.Push(token)
				continue outer
			}
		}
		for _, o := range operators {
			if o == token {
				this.operators.Push(token)
				continue outer
			}
		}
		//check if literal
		_, err := strconv.Atof(token)
		if err == nil {
			this.literals.Push(token)
			continue outer
		}
		//if it reaches here, then it is an identifier
		this.identifiers.Push(token)
	}
}

func (this *Lexicon) Html() string {
	html := "<table cellpadding=\"10\" width=\"700px;\"> \n <tr>"
	data := make([][][2]string, 5)
	vals := [5]vector.StringVector{
		this.keywords,
		this.separators,
		this.operators,
		this.literals,
		this.identifiers,
	}
	//replace with distinct
	for i, v := range vals {
		data[i] = distinct(v)
	}
	lens := make([]int, 5)
	//get the longest values
	for i, d := range data {
		lens[i] = len(d)
	}
	maxLen := func(array []int) int {
		max := array[0]
		for i := 1; i < 5; i++ {
			if max < array[i] {
				max = array[i]
			}
		}
		return max
	}(lens)

	headings := [5]string{"KEYWORDS", "SEPARATORS", "OPERATORS", "LITERALS", "IDENTIFIERS"}
	for _, h := range headings {
		html += fmt.Sprintf("<td style=\"background-color:#CCF;\">%s</td>", h)
	}
	html += "\n</tr>"
	background := false
	for i := 0; i < maxLen; i++ {
		background = func() bool{
			if background {
				return false
			}
			return true
		}()
		back := ""
		if background { back = ` style="background-color:#CCC;"`} else { back = "" }
		html += fmt.Sprintf("\n<tr%s>", back)
		for _, d := range data {
			html += fmt.Sprintf("<td>%s</td>", elemAt(d, i))
		}
		html += "\n</tr>"
	}
	return html + "\n</table>"
}

func distinct(list vector.StringVector) [][2]string {
	if len(list) == 0 {
		return nil
	}
	var keys vector.StringVector
	var vals vector.IntVector
	for _, l := range list {
		index := search(keys, l)
		if index == -1{
			keys.Push(l)
			vals.Push(1)
		}else{
			vals.Set(index, vals.At(index) + 1)
		}
	}
	m := make([][2]string, len(keys))
	for i, k := range keys {
		m[i][0] = k
		m[i][1] = fmt.Sprint(vals.At(i))
	}
	return m
}

func elemAt(list [][2]string, index int) string {
	if len(list) <= index {
		return "&nbsp;"
	}
	return list[index][0] + func() string {
		if list[index][1] == "1"{
			return ""
		}
		return "  &nbsp;- " + list[index][1] + " times"
	}()
}

func search(list []string, key string) int{
	for i, l := range list {
		if l == key {
			return i
		}
	}
	return -1
}
