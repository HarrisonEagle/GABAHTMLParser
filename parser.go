package GABAHTMLParser

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"os"
	"strings"
)

func readLinesFromFile(path string,isEncode bool) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	
	if isEncode{
		dec := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())
		scanner = bufio.NewScanner(dec)
	}
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func readLinesFromURL(url string,isEncode bool) ([]string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return LinesFromReader(resp.Body,isEncode)
}

func LinesFromReader(r io.Reader,isEncode bool) ([]string, error) {
	var lines []string
	dec := r
	if isEncode{
		dec = transform.NewReader(r, japanese.ShiftJIS.NewDecoder())
	}
	buf := make([]byte, 0, 64*1024)
	scanner := bufio.NewScanner(dec)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func GetHTMLfromURL(path string,isEncode bool) *Element{
	lines, err := readLinesFromURL(path,isEncode)
	if err != nil {
		file_lines, file_err := readLinesFromFile(path,isEncode)
		lines = file_lines
		err = file_err
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ParseHTML(lines)
}

func Split(r rune) bool {
	return r == '|' || r == ' '
}

func IsSpecialTag(category string) bool {
	switch category {
	case
		"meta",
		"link",
		"base",
		"input",
		"frame",
		"hr",
		"img":
		return true
	}
	return false
}

func ParseHTML(lines []string) *Element {
	pointer := new(Element)
	pointer.Parent = new(Element)
	pointer.Parent.Child = append(pointer.Parent.Child, pointer)
	isComment := false
	isFillingAttr := false
	TagOpen := false
	stacksize := 0
	temp := 0
	key := ""
	for i, line := range lines {
		temp+=i
		line = strings.Replace(line, "<", "|<", -1)
		line = strings.Replace(line, ">", ">|", -1)
		words := strings.FieldsFunc(line, Split)
		for _, word := range words {
			//time.Sleep(time.Second)
			if TagOpen && (!strings.HasPrefix(word, "<") || !strings.HasPrefix(word, "</") || strings.HasPrefix(word, "<!")) {
				if pointer.Parent.InnerHTML == "" {
					pointer.Parent.InnerHTML = word
				} else {
					pointer.Parent.InnerHTML += (" " + word)
				}
			}
			if  strings.HasPrefix(word, "<!--") {
				isComment = true
				if strings.HasSuffix(word, "-->") {
					isComment = false
				}
			} else if isComment && strings.HasSuffix(word, "-->") {
				isComment = false
			} else if !isComment {
				if strings.HasPrefix(word, "</") || strings.HasSuffix(word, "/>") {
					//Close Tag
					if pointer.Parent != nil {
						if pointer.Parent.InnerHTML == "" {
							pointer.Parent.InnerHTML = pointer.InnerHTML
						} else {
							pointer.Parent.InnerHTML += (" " + pointer.InnerHTML)
						}
						pointer.Parent.InnerHTML += (" " + word)
						pointer = pointer.Parent
					}
					stacksize--
					TagOpen = false
					isFillingAttr = false
				} else if strings.HasPrefix(word, "<") && pointer.Tag != "script" && pointer.Tag != "noscript" && pointer.Tag != "style"{
					TagOpen = true
					isFillingAttr = false
					stacksize++
					var Child *Element = new(Element)
					Child.InnerHTML = ""
					Child.Attr = make(map[string]string)
					Child.Parent = pointer
					pointer.Child = append(pointer.Child, Child)
					if pointer.InnerHTML == "" {
						pointer.InnerHTML = word
					} else {
						pointer.InnerHTML += (" " + word)
					}
					pointer = pointer.Child[len(pointer.Child)-1]
					pointer.Tag = word[1:]
					if strings.HasSuffix(pointer.Tag, ">") {
						pointer.Tag = strings.Replace(pointer.Tag, ">", "", -1)
						TagOpen = false
					}
					/*
					fmt.Print(pointer.Tag+" ")
					fmt.Print(stacksize)
					fmt.Print("  line:")
					fmt.Println(i)
					*/
					if pointer.Tag == "br" || pointer.Tag == "wbr"{
						pointer = pointer.Parent
						stacksize--
						TagOpen = false
						isFillingAttr = false
					}
				} else if !TagOpen {
					if pointer.InnerHTML == "" {
						pointer.InnerHTML = word
					} else {
						pointer.InnerHTML += (" " + word)
					}
				} else {
					if isFillingAttr {
						val := word
						if strings.HasSuffix(word, ">") {
							if len(val) >= 2{
								val = val[:len(val)-2]
							}
							isFillingAttr = false
						} else if strings.HasSuffix(word, "'") || strings.HasSuffix(word, "\"") {
							val = val[:len(val)-1]
							isFillingAttr = false
						}
						pointer.Attr[key] += (" " + val)
					} else {
						isFillingAttr = true
						arr := strings.Split(word, "=\"")
						if len(arr) <= 1 {
							arr = strings.Split(word, "='")
						}
						if len(arr) > 1 {
							key = arr[0]
							val := arr[1]
							if strings.HasSuffix(arr[1], ">") {
								val = val[:len(val)-2]
								isFillingAttr = false
							} else if strings.HasSuffix(arr[1], "'") || strings.HasSuffix(arr[1], "\"") {
								val = val[:len(val)-1]
								isFillingAttr = false
							}
							pointer.Attr[key] = val
						}else{
							isFillingAttr = false
						}
					}
					if strings.HasSuffix(word, ">") && !strings.HasSuffix(word, "/>"){
						if IsSpecialTag(pointer.Tag){
							//Close Tag
							if pointer.Parent != nil {
								pointer = pointer.Parent
							}
							stacksize--
						}
						isFillingAttr = false
						TagOpen = false
					}
				}
			}
		}
	}
	return pointer
}
