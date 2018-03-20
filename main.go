package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"./src"
	"path/filepath"
)

func main() {
	fmt.Println("Hello Gocast!")

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	parser := argparse.NewParser("gocast", "gocast is a program for downloading podcasts from acast.com via rss feed.")
	name := parser.String("n", "name", &argparse.Options{Required: true, Help: "Podcast name e.g. letstalkaboutcarsyo for http://rss.acast.com/letstalkaboutcarsyo"})
	outputPath := parser.String("o", "output", &argparse.Options{
		Required: false,
		Help:     "Specifies where downloaded episodes will be saved to",
		Default:  cwd})
	latest := parser.Flag( "", "latest", &argparse.Options{Help: "Download the latest episode"})
	err = parser.Parse(os.Args)

	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	client, err := gocast.NewAcastClient(fmt.Sprintf("http://rss.acast.com/%s", *name))

	if *latest {
		err = client.DownloadLatestEpisode(*outputPath)
	}

	if err != nil {
		println(err)
	}
}
