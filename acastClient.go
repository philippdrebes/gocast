package gocast

import (
	"github.com/beevik/etree"
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"io"
	"gopkg.in/cheggaaa/pb.v1"
	"path/filepath"
	"regexp"
	"sort"
)

var invalidFileNameRegex = regexp.MustCompile(regexp.QuoteMeta(`[\\/:"*?<>|]+`))

type AcastClient struct {
	url      string
	document *etree.Document
	channel  *etree.Element
	episodes map[int]*Episode
}

type Episode struct {
	title string
	media string
}

func NewAcastClient(url string) (*AcastClient, error) {
	client := &AcastClient{}
	client.url = url

	document, err := client.reload()
	if err != nil {
		return nil, err
	}
	client.document = document

	root := client.document.SelectElement("rss")
	client.channel = root.SelectElement("channel")

	fmt.Println(client.channel.SelectElement("title").Text())

	client.episodes = *client.GetAllEpisodes()

	return client, nil
}

func (c AcastClient) reload() (*etree.Document, error) {
	response, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromString(string(content)); err != nil {
		panic(err)
	}

	return doc, nil
}

func (c AcastClient) GetAllEpisodes() *map[int]*Episode {
	if c.channel == nil {
		return nil
	}

	episodes := make(map[int]*Episode)

	for index, episode := range c.channel.SelectElements("item") {
		title := episode.SelectElement("title")
		enclosure := episode.SelectElement("enclosure")
		media := enclosure.SelectAttrValue("url", "")
		if len(media) == 0 {
			fmt.Printf("\nCould not find media link for %s", title)
		} else {
			episodes[index] = &Episode{title: title.Text(), media: media}
		}
	}

	return &episodes
}

func (c AcastClient) ListAllEpisodes() {
	if c.episodes == nil {
		return
	}

	var keys []int
	for k := range c.episodes {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		fmt.Printf("%d) %s\n", key, c.episodes[key].title)
	}
}

func (c AcastClient) DownloadAllEpisodes(outputPath string) error {
	if c.channel == nil {
		return nil
	}

	fmt.Printf("Downloading all episodes")

	episodeCount := len(c.episodes)
	if c.episodes == nil || episodeCount == 0 {
		fmt.Println("Could not find episodes")
	}

	for index, episode := range c.episodes {
		filename := invalidFileNameRegex.ReplaceAllString(episode.title, "")

		fmt.Printf("Downloading [%d/%d]: %s\n", index+1, episodeCount, episode.title)
		err := c.download(episode.media, filepath.Join(outputPath, fmt.Sprintf("%s.mp3", filename)))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c AcastClient) DownloadEpisode(index int, outputPath string) error {
	if c.channel == nil {
		return nil
	}

	episode := c.episodes[index]
	if len(episode.media) == 0 {
		fmt.Printf("\nError downloading %s. Could not find media link.", episode.title)
		return nil
	}

	filename := invalidFileNameRegex.ReplaceAllString(episode.title, "")
	fmt.Printf("Downloading episode: %s\n", episode.title)

	outputPath = filepath.Join(outputPath, fmt.Sprintf("%s.mp3", filename))
	c.download(episode.media, outputPath)

	return nil
}

func (c AcastClient) download(url string, output string) (err error) {
	// create the file
	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	// create and start bar
	bar := pb.New(int(response.ContentLength)).SetUnits(pb.U_BYTES)
	bar.Start()

	// create proxy reader
	reader := bar.NewProxyReader(response.Body)

	// check server response
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", response.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}
