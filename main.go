package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/urfave/cli/v2"
)

var (
	blockheader string
	filename    string
	matchstr    string
	delim       bool
	nohighlight bool
	ignorecase  bool
	blocks      []string
)

func main() {
	app := &cli.App{
		Name:                   "logrep",
		Usage:                  "log grep",
		UsageText:              "logrep [options] <matchstr> [filename]",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "blockheader",
				Aliases: []string{"b"},
				Value:   "^[0-9]{4} [0-9]{2}:",
				Usage:   "block header regexp",
				// Required: true,
			},
			&cli.BoolFlag{
				Name:  "delim",
				Usage: "print line delim",
			},
			&cli.BoolFlag{
				Name:  "nohighlight",
				Usage: "No Highlight",
			},
			&cli.BoolFlag{
				Name:    "ignorecase",
				Aliases: []string{"i"},
				Usage:   "ignore case",
			},
		},
		Action: func(c *cli.Context) error {
			matchstr = c.Args().Get(0)
			filename = c.Args().Get(1)

			if matchstr == "" {
				cli.ShowAppHelp(c)
				os.Exit(1)
			}

			blockheader = c.String("blockheader")
			delim = c.Bool("delim")
			nohighlight = c.Bool("nohighlight")
			ignorecase = c.Bool("ignorecase")

			fmt.Printf(`
	match:        "%s"
	file:         "%s"
	blockheader:  "%s"
	delim:        %v
	nohighlight:  %v
	ignorecase:   %v
	
`,
				matchstr, filename, blockheader, delim, nohighlight, ignorecase)

			DoMatch()

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func DoMatch() {
	if ignorecase {
		matchstr = "(?i)" + matchstr
	}
	reg, err := regexp.Compile(matchstr)
	if err != nil {
		log.Fatalf("matchstr regexp compile error: %v\n", err)
	}

	regblock, err := regexp.Compile(blockheader)
	if err != nil {
		log.Fatalf("blockheader regexp compile error: %v\n", err)
	}

	var f *os.File

	if filename != "" {
		f, err = os.Open(filename)
		if err != nil {
			log.Fatalf(`open file "%s" error: %v`+"\n", filename, err)
		}
		defer f.Close()
	} else {
		f = os.Stdin
	}

	found := false
	totalnum := 0

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		if regblock.FindStringIndex(line) != nil {
			if found {
				//matching block end
				for _, v := range blocks {
					fmt.Println(v)
				}
				totalnum++

				if delim {
					fmt.Printf("\n----------\n\n")
				}

				found = false
			}

			blocks = nil
			blocks = append(blocks, line)
		} else {
			blocks = append(blocks, line)
		}

		if reg.FindStringIndex(line) != nil {
			// log.Println("found!")
			found = true

			if !nohighlight && blocks != nil {
				blocks[len(blocks)-1] = reg.ReplaceAllStringFunc(line, func(src string) string {
					return Yellow(src)
				})
			}
		}
	}

	if found {
		//matching block end
		for _, v := range blocks {
			fmt.Println(v)
		}
		totalnum++
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
