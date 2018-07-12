package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"path/filepath"
	"github.com/philippdrebes/gocast"
)

func main() {
	fmt.Println("Hello Gocast!")

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}

	parser := argparse.NewParser("gocast", "gocast is a program for downloading podcasts from acast.com via rss feed.")
	name := parser.String("n", "name", &argparse.Options{
		Required: true,
		Help: "Podcast name e.g. letstalkaboutcarsyo for http://rss.acast.com/letstalkaboutcarsyo"})
	outputPath := parser.String("o", "output", &argparse.Options{
		Required: false,
		Help:     "Specifies where downloaded episodes will be saved to",
		Default:  cwd})
	list := parser.Flag("l", "list", &argparse.Options{Help: "List all episodes"})
	index := parser.Int("i", "index", &argparse.Options{
		Help: "Download a single episode via index. Run the 'list' command in order to get the index of the desired episode",
		Default: -1})
	all := parser.Flag("", "all", &argparse.Options{Help: "Download all episodes"})
	latest := parser.Flag("", "latest", &argparse.Options{Help: "Download the latest episode"})
	err = parser.Parse(os.Args)

	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	client, err := gocast.NewAcastClient(fmt.Sprintf("http://rss.acast.com/%s", *name))

	if *list {
		client.ListAllEpisodes()
	}
	if *index != -1 {
		client.DownloadEpisode(*index, *outputPath)
	}
	if *latest {
		err = client.DownloadEpisode(0, *outputPath)
	}
	if *all {
		err = client.DownloadAllEpisodes(*outputPath)
	}

	if err != nil {
		println(err)
	}
}
