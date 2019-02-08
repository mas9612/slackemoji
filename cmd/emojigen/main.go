package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/mas9612/slackemoji/pkg/slackemoji"
)

var (
	color   string
	font    string
	public  bool
	text    string
	outfile string
	endian  string
)

func main() {
	flag.StringVar(&color, "color", "", "Font color code omitted precede # character. Default: EC71A1FF")
	flag.StringVar(&font, "font", "", "Font family. Default: notosans-mono-bold. Valid family is one of [notosans-mono-bold, mplus-1p-black, rounded-x-mplus-1p-black, ipamjm, LinLibertine_RBah, aoyagireisyoshimo].")
	flag.BoolVar(&public, "public", false, "Whether new emoji will be make public. Default: false.")
	flag.StringVar(&text, "text", "", "Emoji text. Required. If you want to generate with multiline text, separate each line with comma character (,).")
	flag.StringVar(&outfile, "out", "", "Output filename. Required.")
	flag.Parse()

	if text == "" {
		fmt.Fprintf(os.Stderr, "-text is required")
		os.Exit(1)
	}

	var options []slackemoji.EmojiOption
	if color != "" {
		options = append(options, slackemoji.Color(color))
	}
	if font != "" {
		switch font {
		case "notosans-mono-bold", "mplus-1p-black", "rounded-x-mplus-1p-black", "ipamjm", "LinLibertine_RBah", "aoyagireisyoshimo":
			options = append(options, slackemoji.Font(font))
		default:
			fmt.Fprintf(os.Stderr, "Invalid font family '%s'", font)
			os.Exit(1)
		}
	}
	if public {
		options = append(options, slackemoji.Public(public))
	}
	if outfile == "" {
		fmt.Fprintf(os.Stderr, "-out is required.")
		os.Exit(1)
	}

	emoji, err := slackemoji.GenerateEmoji(text, options...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate emoji: %v", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open new file '%s'", outfile)
		os.Exit(1)
	}
	defer file.Close()
	binary.Write(file, binary.LittleEndian, emoji)
}
