package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"html/template"
	"strings"
)

type DevNullWriter struct{}

func (w *DevNullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

var devnull = &DevNullWriter{}
var sb = &strings.Builder{}

func FuzzXSSpoc(data []byte) int {
	t, err := template.New("fuzz").Parse(`<p href={{.}}>`)
	if err != nil {
		return 0
	}
	sb.Reset()
	err = t.Execute(sb, data)
	if err != nil {
		return 0
	}
	_, err = html.Parse(strings.NewReader(sb.String()))
	if err != nil {
		return 0
	}
	sb.Reset()
	err = t.Execute(sb, string(""))
	if err != nil {
		panic("XSS")
	}
	_, err = html.Parse(strings.NewReader(sb.String()))
	fmt.Printf("lol %s\n", sb.String())
	if err != nil {
		panic("XSS")
	}
	return 1
}

func FuzzXSS(data []byte) int {
	sep := bytes.IndexByte(data, 0)
	if sep < 0 || sep+1 > len(data) {
		return 0
	}
	//<p href={{.}}>
	t, err := template.New("fuzz").Parse(string(data[:sep]))
	if err != nil {
		return 0
	}
	sb.Reset()
	err = t.Execute(sb, "A")
	if err != nil {
		return 0
	}
	_, err = html.Parse(strings.NewReader(sb.String()))
	if err != nil {
		panic("XSS1")
	}
	sb.Reset()
	err = t.Execute(sb, string(data[sep+1:]))
	if err != nil {
		panic("XSS Execute")
	}
	_, err = html.Parse(strings.NewReader(sb.String()))
	if err != nil {
		panic("XSS Parse")
	}

	return 1
}
