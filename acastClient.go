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
	"strings"
)

type AcastClient struct {
	url      string
	document *etree.Document
	channel  *etree.Element
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

func (c AcastClient) DownloadAllEpisodes(outputPath string) error {
	if c.channel == nil {
		return nil
	}

	fmt.Printf("Downloading all episodes")

	episodeCount := len(c.channel.SelectElements("item"))
	for index, episode := range c.channel.SelectElements("item") {
		title := episode.SelectElement("title")
		enclosure := episode.SelectElement("enclosure")

		media := enclosure.SelectAttrValue("url", "")
		if len(media) == 0 {
			fmt.Printf("\nError downloading %s. Could not find media link.", title)
			return nil
		}
		fmt.Printf("Downloading [%d/%d]: %s\n", index+1, episodeCount, title.Text())
		err := c.download(media, filepath.Join(outputPath, fmt.Sprintf("%s.mp3", title.Text())))

		if err != nil {
			return err
		}
	}

	return nil
}

func (c AcastClient) DownloadLatestEpisode(outputPath string) error {
	if c.channel == nil {
		return nil
	}

	//author := c.channel.SelectElement("itunes:author")

	for _, episode := range c.channel.SelectElements("item")[:1] {
		title := episode.SelectElement("title")
		enclosure := episode.SelectElement("enclosure")

		media := enclosure.SelectAttrValue("url", "")
		if len(media) == 0 {
			fmt.Printf("\nError downloading %s. Could not find media link.", title)
			return nil
		}

		titleText := title.Text()
		titleText = strings.Replace(titleText, ":", "", -1)
		fmt.Printf("Downloading latest episode: %s\n", title.Text())
		outputPath = filepath.Join(outputPath, fmt.Sprintf("%s.mp3", titleText))

		c.download(media, outputPath)

		// set ID3 Tags
		//mp3File, _ := id3.Open(outputPath)
		//mp3File.SetTitle(title.Text())
		//
		//if author != nil {
		//	mp3File.SetArtist(author.Text())
		//}

		//desc := episode.FindElement("itunes:summary")
		//if desc != nil {
		//	ft := v2.V23FrameTypeMap["TIT3"]
		//	descText := desc.Text()
		//	descFrame := v2.NewTextFrame(ft, descText)
		//	mp3File.AddFrames(descFrame)
		//}

		//picture := episode.FindElement("itunes:image")
		//if picture != nil {
		//	pictureHref := picture.SelectAttrValue("href", "")
		//	if len(pictureHref) > 0 {
		//
		//		response, err := http.Get(pictureHref)
		//		if err != nil {
		//			return err
		//		}
		//
		//		defer response.Body.Close()
		//
		//		ft := v2.V23FrameTypeMap["APIC"]
		//
		//		buf, err := ioutil.ReadAll(response.Body)
		//		pictureFrame := v2.NewDataFrame(ft, buf)
		//		mp3File.AddFrames(pictureFrame)
		//	}
		//}

		//mp3File.Close()
	}

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
