package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// Org mode represantation of Node
func (nodes Nodes) orgFormatOriginal() string {

	value := ""
	header := 0
	table := 0
	list := 0
	listValue := []int{}

	for _, node := range nodes {
		switch node.Token.Type {

		case html.SelfClosingTagToken:
			switch node.Token.Data {
			case "en-media":
				base, file := mimeFiles(node.Token)
				value += "[[./" + base + "/"
				value += file + "]]"

			case "en-todo":

				switch getAttr("checked", node.Token) {
				case "true":
					value += "\n- [X] "
				case "false":
					value += "\n- [ ] "
				}
			}
		case html.StartTagToken:
			switch node.Token.Data {

			case "a":
				// We do not want links in the header
				if header == 0 {
					value += "[[" + getAttr("href", node.Token) + "]["
				}
			case "p":
				value += "\n"
			case "u":
				value += "_"
			case "i":
				value += "/"
			case "b", "strong", "em":
				value += "*"
			case "del":
				value += "+"
			case "h1":
				value += "\n** "
				header++
			case "h2":
				value += "\n*** "
				header++
			case "h3":
				value += "\n**** "
				header++
			case "h4":
				value += "\n***** "
				header++
			case "h5":
				value += "\n****** "
				header++
			case "h6":
				value += "\n******* "
				header++

				// These tags are ignored
			case "div", "span", "tr", "tbody", "abbr", "th", "thead", "ins", "img":
				break
			case "sup", "small", "br", "dl", "dd", "dt", "font", "colgroup", "cite":
				break
			case "address", "s", "map", "area", "center":
				break

			case "hr":
				value += "\n------\n"
			case "en-media":
				base, file := mimeFiles(node.Token)
				value += "[[./" + base + "/"
				value += file + "]]"
			case "table":
				table++
			case "td":
				value += "|"
			case "ol":
				list++
				listValue = append(listValue, 1)
			case "ul":
				list++
				listValue = append(listValue, 0)
			case "li":
				value += "\n"
				for i := 0; i <= list; i++ {
					value += "  "
				}
				if list > 0 {
					switch listValue[list-1] {
					case 0:
						value += "- "
					default:
						value += fmt.Sprintf("%d.", listValue[list-1])
						listValue[list-1] = listValue[list-1] + 1
					}
				}
			case "code":
				value += "~"
			case "pre":
				value += "\n#+BEGIN_SRC\n"
			case "blockquote":
				value += "\n#+BEGIN_QUOTE\n"

			default:
				println(node.Token.Data)
			}

		case html.EndTagToken:
			switch node.Token.Data {
			case "u":
				value += "_"
			case "i":
				value += "/"
			case "b", "strong", "em":
				value += "*"
			case "del":
				value += "+"
			case "a":
				if header == 0 {
					value += "]]"
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				header--
			case "table":
				table--
			case "tr":
				value += "|\n"
			case "ol", "ul":
				list--
				listValue = listValue[:len(listValue)-1]
			case "code":
				value += "~"
			case "pre":
				value += "\n#+END_SRC\n"
			case "blockquote":
				value += "\n#+END_QUOTE\n"

			}
		}
		value += node.Text
	}
	return value
}

// Org mode representation of Node
func (nodes Nodes) orgFormatStringsBuilder() string {

	var value strings.Builder
	header := 0
	table := 0
	list := 0
	listValue := []int{}

	for _, node := range nodes {
		switch node.Token.Type {

		case html.SelfClosingTagToken:
			switch node.Token.Data {
			case "en-media":
				base, file := mimeFiles(node.Token)
				value.WriteString("[[./" + base + "/")
				value.WriteString(file + "]]")

			case "en-todo":

				switch getAttr("checked", node.Token) {
				case "true":
					value.WriteString("\n- [X] ")
				case "false":
					value.WriteString("\n- [ ] ")
				}
			}
		case html.StartTagToken:
			switch node.Token.Data {

			case "a":
				// We do not want links in the header
				if header == 0 {
					value.WriteString("[[" + getAttr("href", node.Token) + "][")
				}
			case "p":
				value.WriteString("\n")
			case "u":
				value.WriteString("_")
			case "i":
				value.WriteString("/")
			case "b", "strong", "em":
				value.WriteString("*")
			case "del":
				value.WriteString("+")
			case "h1":
				value.WriteString("\n** ")
				header++
			case "h2":
				value.WriteString("\n*** ")
				header++
			case "h3":
				value.WriteString("\n**** ")
				header++
			case "h4":
				value.WriteString("\n***** ")
				header++
			case "h5":
				value.WriteString("\n****** ")
				header++
			case "h6":
				value.WriteString("\n******* ")
				header++

				// These tags are ignored
			case "div", "span", "tr", "tbody", "abbr", "th", "thead", "ins", "img":
				break
			case "sup", "sub", "small", "br", "dl", "dd", "dt", "font", "colgroup", "cite":
				break
			case "address", "s", "map", "area", "center", "q":
				break

			case "hr":
				value.WriteString("\n------\n")
			case "en-media":
				base, file := mimeFiles(node.Token)
				value.WriteString("[[./" + base + "/")
				value.WriteString(file + "]]")
			case "table":
				table++
			case "td":
				value.WriteString("|")
			case "ol":
				list++
				listValue = append(listValue, 1)
			case "ul":
				list++
				listValue = append(listValue, 0)
			case "li":
				value.WriteString("\n")
				for i := 0; i <= list; i++ {
					value.WriteString("  ")
				}
				if list > 0 {
					switch listValue[list-1] {
					case 0:
						value.WriteString("- ")
					default:
						value.WriteString(fmt.Sprintf("%d.", listValue[list-1]))
						listValue[list-1] = listValue[list-1] + 1
					}
				}
			case "code":
				value.WriteString("~")
			case "pre":
				value.WriteString("\n#+BEGIN_SRC\n")
			case "blockquote":
				value.WriteString("\n#+BEGIN_QUOTE\n")

			default:
				fmt.Println("skip token: " + node.Token.Data)
			}

		case html.EndTagToken:
			switch node.Token.Data {
			case "u":
				value.WriteString("_")
			case "i":
				value.WriteString("/")
			case "b", "strong", "em":
				value.WriteString("*")
			case "del":
				value.WriteString("+")
			case "a":
				if header == 0 {
					value.WriteString("]]")
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				header--
			case "table":
				table--
			case "tr":
				value.WriteString("|\n")
			case "ol", "ul":
				list--
				listValue = listValue[:len(listValue)-1]
			case "code":
				value.WriteString("~")
			case "pre":
				value.WriteString("\n#+END_SRC\n")
			case "blockquote":
				value.WriteString("\n#+END_QUOTE\n")

			}
		}
		value.WriteString(node.Text)
	}
	return value.String()
}

// Org mode representation of Node
func (nodes Nodes) orgFormatStringBuffer() string {

	value := &bytes.Buffer{}

	header := 0
	table := 0
	list := 0
	listValue := []int{}

	for _, node := range nodes {
		switch node.Token.Type {

		case html.SelfClosingTagToken:
			switch node.Token.Data {
			case "en-media":
				base, file := mimeFiles(node.Token)
				value.WriteString("[[./" + base + "/")
				value.WriteString(file + "]]")

			case "en-todo":

				switch getAttr("checked", node.Token) {
				case "true":
					value.WriteString("\n- [X] ")
				case "false":
					value.WriteString("\n- [ ] ")
				}
			}
		case html.StartTagToken:
			switch node.Token.Data {

			case "a":
				// We do not want links in the header
				if header == 0 {
					value.WriteString("[[" + getAttr("href", node.Token) + "][")
				}
			case "p":
				value.WriteString("\n")
			case "u":
				value.WriteString("_")
			case "i":
				value.WriteString("/")
			case "b", "strong", "em":
				value.WriteString("*")
			case "del":
				value.WriteString("+")
			case "h1":
				value.WriteString("\n** ")
				header++
			case "h2":
				value.WriteString("\n*** ")
				header++
			case "h3":
				value.WriteString("\n**** ")
				header++
			case "h4":
				value.WriteString("\n***** ")
				header++
			case "h5":
				value.WriteString("\n****** ")
				header++
			case "h6":
				value.WriteString("\n******* ")
				header++

				// These tags are ignored
			case "div", "span", "tr", "tbody", "abbr", "th", "thead", "ins", "img":
				break
			case "sup", "sub", "small", "br", "dl", "dd", "dt", "font", "colgroup", "cite":
				break
			case "address", "s", "map", "area", "center", "q":
				break

			case "hr":
				value.WriteString("\n------\n")
			case "en-media":
				base, file := mimeFiles(node.Token)
				value.WriteString("[[./" + base + "/")
				value.WriteString(file + "]]")
			case "table":
				table++
			case "td":
				value.WriteString("|")
			case "ol":
				list++
				listValue = append(listValue, 1)
			case "ul":
				list++
				listValue = append(listValue, 0)
			case "li":
				value.WriteString("\n")
				for i := 0; i <= list; i++ {
					value.WriteString("  ")
				}
				if list > 0 {
					switch listValue[list-1] {
					case 0:
						value.WriteString("- ")
					default:
						value.WriteString(fmt.Sprintf("%d.", listValue[list-1]))
						listValue[list-1] = listValue[list-1] + 1
					}
				}
			case "code":
				value.WriteString("~")
			case "pre":
				value.WriteString("\n#+BEGIN_SRC\n")
			case "blockquote":
				value.WriteString("\n#+BEGIN_QUOTE\n")

			default:
				fmt.Println("skip token: " + node.Token.Data)
			}

		case html.EndTagToken:
			switch node.Token.Data {
			case "u":
				value.WriteString("_")
			case "i":
				value.WriteString("/")
			case "b", "strong", "em":
				value.WriteString("*")
			case "del":
				value.WriteString("+")
			case "a":
				if header == 0 {
					value.WriteString("]]")
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				header--
			case "table":
				table--
			case "tr":
				value.WriteString("|\n")
			case "ol", "ul":
				list--
				listValue = listValue[:len(listValue)-1]
			case "code":
				value.WriteString("~")
			case "pre":
				value.WriteString("\n#+END_SRC\n")
			case "blockquote":
				value.WriteString("\n#+END_QUOTE\n")

			}
		}
		value.WriteString(node.Text)
	}
	return value.String()
}

var (
	content, _ = ioutil.ReadFile("sample_content.html")
	reader     = bytes.NewReader(content)

	nodes = parseHTML(reader)
)

func BenchmarkOrgFormatOriginal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = nodes.orgFormatOriginal()
	}
}
func BenchmarkOrgConcatStringBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = nodes.orgFormatStringBuffer()
	}
}

func BenchmarkOrgConcatStringBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = nodes.orgFormatStringsBuilder()
	}
}
