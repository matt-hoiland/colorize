package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"regexp"

	"github.com/fatih/color"
)

type matcher struct {
	re    *regexp.Regexp
	print func(string, ...interface{})
}

var matchers = []matcher{
	{re: regexp.MustCompile("\"level\":\"(error|ERROR)\""), print: color.Red},
	{re: regexp.MustCompile("\"level\":\"(warning|WARNING)\""), print: color.Yellow},
	{re: regexp.MustCompile("\"level\":\"(info|INFO)\""), print: color.Cyan},
	{re: regexp.MustCompile("\"level\":\"(debug|DEBUG)\""), print: color.Green},
	{re: regexp.MustCompile(".*"), print: func(s string, a ...interface{}) { fmt.Println(s) }}, // Default
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if info.Mode()&fs.ModeNamedPipe == 0 {
		fmt.Println("This program should be used as a pipe.")
		fmt.Println("    ex: `echo 'hello' | colorize")
		return
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		for _, m := range matchers {
			if m.re.Match(line) {
				m.print(string(line))
				break
			}
		}
	}
}
