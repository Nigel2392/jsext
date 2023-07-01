package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Example go file:
// The following file defines a new function called SvgFunc, which returns a new svg element.
// Since html is case insensitive, the function name is converted to lowercase.
// We can change this by specifying the upper attribute.
//
// When not specifying a newfunc, the function to create a new element will be jse.NewElement
//
// Imports can be specified with the imports attribute.
// This attribute should be a list of items, separated by a semicolon ;
//
// This example does not nescessarily work, it is only to show what can be done.
//
//		//go:generate jsextgen $GOFILE
//		/*
//		<jsextgen>
//			<svgFunc upper="s;f" newfunc="svg.New" imports="syscall/js;net/http"
//				<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-0-square-fill" viewBox="0 0 16 16">
//				  <path d="M8 4.951c-1.008 0-1.629 1.09-1.629 2.895v.31c0 1.81.627 2.895 1.629 2.895s1.623-1.09 1.623-2.895v-.31c0-1.8-.621-2.895-1.623-2.895Z"/>
//				  <path d="M2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2H2Zm5.988 12.158c-1.851 0-2.941-1.57-2.941-3.99V7.84c0-2.408 1.101-3.996 2.965-3.996 1.857 0 2.935 1.57 2.935 3.996v.328c0 2.408-1.101 3.99-2.959 3.99Z"/>
//				</svg>
//			</svgFunc>
//		</jsextgen>
//		*/

const (
	OPENING_TAG = "<jsextgen>"
	CLOSING_TAG = "</jsextgen>"

	PACKAGE = "jsextgen"
)

func main() {
	if len(os.Args) < 2 {
		panic("No file specified.")
	}

	fileName := os.Args[1]

	var flagParser = flag.NewFlagSet(PACKAGE, flag.ExitOnError)

	var outFile = flagParser.String("o", PACKAGE+".go", "Output file")
	var packageName = flagParser.String("p", PACKAGE, "Package name")

	flagParser.Parse(os.Args[2:])

	if *outFile == "" {
		panic("No output file specified.")
	}

	if *packageName == "" {
		panic("No package name specified.")
	}

	// parse the go file, get all comments and put the comments into a slice for later use
	var goFile, err = parser.ParseFile(token.NewFileSet(), fileName, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generating file:", *outFile, "from file:", fileName)

	var f = NewMultiBufferedWriter(func() LenStringerWriter { return &strings.Builder{} })

	var imports = make(map[string]struct{})
	var funcs = make(map[string]struct{})

	for _, comment := range goFile.Comments {
		var comment = comment.Text()

		var split = strings.Split(comment, CLOSING_TAG)
	htmlLoop:
		for _, functionHTML := range split {

			if strings.TrimSpace(functionHTML) == "" {
				continue
			}

			f.New()

			var htmlString = strings.TrimSpace(functionHTML)
			htmlString = strings.TrimPrefix(htmlString, OPENING_TAG)
			strings.TrimSuffix(htmlString, CLOSING_TAG)
			if strings.TrimSpace(htmlString) == "" {
				continue
			}

			var doc, err = html.Parse(strings.NewReader(htmlString))
			if err != nil {
				fmt.Println("Error parsing html:", err)
				continue
			}

			var n *html.Node
		elemLoop1:
			for c := doc.FirstChild.FirstChild.NextSibling.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode {
					n = c
					break elemLoop1
				}
			}
			if n == nil {
				continue htmlLoop
			}

			var funcName = n.Data
			var defaultNewFunc = "jse.NewElement"

			if strings.TrimSpace(funcName) == "" {
				panic("No data in root node (function not defined)")
			}

		attrLoop:
			for _, a := range n.Attr {
				switch a.Key {
				case "upper":
					var parts = strings.Split(a.Val, ";")
					for _, p := range parts {
						funcName = strings.Replace(funcName, p, strings.ToUpper(p), 1)
					}
				case "newfunc":
					defaultNewFunc = strings.TrimSpace(a.Val)
					break attrLoop
				case "imports":
					var parts = strings.Split(a.Val, ";")
					for _, p := range parts {
						imports[strings.TrimSpace(p)] = struct{}{}
					}
				}
			}

			if _, ok := funcs[funcName]; ok {
				panic("Function already defined")
			}

			f.WriteString("func ")
			f.WriteString(funcName)
			f.WriteString("() *jse.Element {\n")

			if n.FirstChild == nil {
				panic("no root element")
			}

			var currentElem *html.Node
		elemLoop2:
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && strings.TrimSpace(c.Data) != "" {
					currentElem = c
					break elemLoop2
				}
			}

			if currentElem == nil {
				panic("no root element")
			}

			f.WriteString("\tvar e = ")
			f.WriteString(defaultNewFunc)
			f.WriteString("(")
			f.WriteString(strconv.Quote(currentElem.Data))
			f.WriteString(")\n")

			for _, a := range currentElem.Attr {
				f.WriteString("\te.SetAttr(\"")
				f.WriteString(a.Key)
				f.WriteString("\", \"")
				f.WriteString(a.Val)
				f.WriteString("\")\n")
			}
			for c := currentElem.FirstChild; c != nil; c = c.NextSibling {
				parseChildren("e", c, f)
			}
			f.WriteString("\treturn e\n")
			f.WriteString("}\n\n")

			fmt.Println("Found function:", funcName)

		}
	}

	os.MkdirAll(*packageName, 0777)
	os.Chdir(*packageName)
	file, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	file.WriteString("package ")
	file.WriteString(*packageName)
	file.WriteString("\n\n")

	file.WriteString("import (\n")
	for i := range imports {
		file.WriteString("\t\"")
		file.WriteString(i)
		file.WriteString("\"\n")
	}
	file.WriteString("\t\"github.com/Nigel2392/jsext/v2/jse\"\n")
	file.WriteString(")\n\n")

	_, err = file.WriteString(f.String())
	if err != nil {
		panic(err)
	}
}

type LenStringerWriter interface {
	Len() int
	String() string
	io.StringWriter
}

type multiBufferedWriter[T LenStringerWriter] struct {
	activeBuf *T
	buffers   []T
	len       int
	newF      func() T
}

func NewMultiBufferedWriter[T LenStringerWriter](newFunc func() T) *multiBufferedWriter[T] {
	var newBuf = newFunc()
	return &multiBufferedWriter[T]{
		activeBuf: &newBuf,
		newF:      newFunc,
	}
}

func (m *multiBufferedWriter[T]) New() {
	var b = m.newF()
	if m.activeBuf == nil {
		m.activeBuf = &b
		return
	}
	m.buffers = append(m.buffers, *m.activeBuf)
	m.activeBuf = &b
}

func (m *multiBufferedWriter[T]) Len() int {
	return m.len
}

func (m *multiBufferedWriter[T]) WriteString(str string) (int, error) {
	m.len += len(str)
	return (*m.activeBuf).WriteString(str)
}

func (m *multiBufferedWriter[T]) String() string {
	var b strings.Builder
	for _, buf := range m.buffers {
		b.WriteString(buf.String())
	}
	b.WriteString((*m.activeBuf).String())
	return b.String()
}

func parseChildren(parent string, n *html.Node, b LenStringerWriter) {
	if strings.TrimSpace(n.Data) == "" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseChildren(parent, c, b)
		}
		return
	}
	switch n.Type {
	case html.TextNode:
		b.WriteString("\t")
		b.WriteString(parent)
		b.WriteString(".InnerText(")
		b.WriteString(strconv.Quote(n.Data))
		b.WriteString(")\n")
		b.WriteString("\n")
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseChildren(parent, c, b)
		}

	case html.ElementNode:
		if n.FirstChild != nil {
			b.WriteString("\n")
		}
		var elem = fmt.Sprintf("%s%s%d", parent, n.Data, b.Len())
		if n.FirstChild != nil || len(n.Attr) > 0 {
			b.WriteString("\tvar ")
			b.WriteString(elem)
			b.WriteString(" = ")
		} else {
			b.WriteString("\t")
		}
		b.WriteString(parent)
		b.WriteString(".NewElement(\"")
		b.WriteString(n.Data)
		b.WriteString("\")\n")
		if len(n.Attr) > 0 {
			for _, a := range n.Attr {
				b.WriteString("\t")
				b.WriteString(elem)
				b.WriteString(".SetAttr(\"")
				b.WriteString(a.Key)
				b.WriteString("\", \"")
				b.WriteString(a.Val)
				b.WriteString("\")\n")
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parseChildren(elem, c, b)
		}
	}
}
