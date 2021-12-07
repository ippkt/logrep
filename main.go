package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	flags "github.com/jessevdk/go-flags"
)

var (
	blockheading string
	filename     string
	matchstr     string
	delim        bool
	nohighlight  bool
	ignorecase   bool
	blocks       []string
)

var opts struct {
	Blockheading string `short:"b" long:"blockheading" default:"^[0-9]{4} [0-9]{2}:" description:"block heading regexp"`
	Delim        bool   `long:"delim" description:"delim"`
	NoHighlight  bool   `long:"nohighlight" description:"No Highlight"`
	IgnoreCase   bool   `short:"i" long:"ignorecase" description:"ignore case"`
	Args         struct {
		Matchstr string
		Filename string
	} `positional-args:"yes" required:"yes"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	blockheading = opts.Blockheading
	matchstr = opts.Args.Matchstr
	filename = opts.Args.Filename
	delim = opts.Delim
	nohighlight = opts.NoHighlight
	ignorecase = opts.IgnoreCase

	fmt.Printf(`
match:        "%s"
file:         "%s"
blockheading: "%s"
delim:        %v
nohighlight:  %v
ignorecase:   %v

`,
		matchstr, filename, blockheading, delim, nohighlight, ignorecase)

	DoMatch()
}

func DoMatch() {
	if ignorecase {
		matchstr = "(?i)" + matchstr
	}
	reg, err := regexp.Compile(matchstr)
	if err != nil {
		log.Panicf("matchstr regexp error: %v\n", err)
	}

	regblock, err := regexp.Compile(blockheading)
	if err != nil {
		log.Panicf("blockheading regexp error: %v\n", err)
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Panicf(`open file "%s" error: %v`+"\n", filename, err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	found := false
	totalnum := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("ReadString err: %v\n", err)
			}

			if found {
				//matching block end
				for _, v := range blocks {
					fmt.Printf("%s", v)
				}
				totalnum++
			}

			break
		}

		if regblock.FindStringIndex(line) != nil {
			if found {
				//matching block end
				for _, v := range blocks {
					fmt.Printf("%s", v)
				}
				totalnum++

				if delim {
					fmt.Printf("\n----------\n\n")
				}

				found = false
			}

			blocks = make([]string, 0)
			blocks = append(blocks, line)
		} else {
			if len(blocks) > 0 {
				blocks = append(blocks, line)
			}
		}

		if reg.FindStringIndex(line) != nil {
			// log.Println("found!")
			found = true

			if !nohighlight && len(blocks) > 0 {
				blocks[len(blocks)-1] = reg.ReplaceAllStringFunc(line, func(src string) string {
					return Yellow(src)
				})
			}
		}
	}

	fmt.Printf("\n==== total %d matches ====\n", totalnum)
}

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) string {
	return textColor(textBlack, str)
}

func Red(str string) string {
	return textColor(textRed, str)
}
func Yellow(str string) string {
	return textColor(textYellow, str)
}
func Green(str string) string {
	return textColor(textGreen, str)
}
func Cyan(str string) string {
	return textColor(textCyan, str)
}
func Blue(str string) string {
	return textColor(textBlue, str)
}
func Purple(str string) string {
	return textColor(textPurple, str)
}
func White(str string) string {
	return textColor(textWhite, str)
}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
