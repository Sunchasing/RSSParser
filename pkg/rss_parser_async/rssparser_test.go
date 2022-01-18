package rss_reader

import (
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var mockXml = func() string {
	xmlFile, _ := os.Getwd()
	xmlFile = filepath.Dir(xmlFile)
	xmlFile = filepath.Dir(xmlFile) + "\\mocks\\mockedXML.xml"
	return xmlFile
}()

type mockHttpClient struct {
	httpInterface
}

func setupParser(maxLen int) (p *rssParserSync, wgp sync.WaitGroup) {
	out := make(chan []RssItem, maxLen)

	p = &rssParserSync{
		wg:  &wgp,
		out: out,
		hr:  &mockHttpClient{},
	}
	return
}

func (r *mockHttpClient) get(url string) (responseBytes []byte, err error) {
	// Replacing a http GET with our known file in order to have a reproducible state
	f, err := os.Open(url)
	if err != nil {
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			return
		}
	}()

	return ioutil.ReadAll(f)
}

func Test_parseToChannel(t *testing.T) {
	urls := []string{mockXml}
	p, wg := setupParser(len(urls))
	readBytes, _ := p.hr.get(urls[0])
	// 117823 is the length of the provided xml file
	require.Equal(t, 117823, len(readBytes))

	require.Equal(t, &wg, p.wg)

	p.wg.Add(1)

	err := p.parseToChannel(urls[0])
	p.wg.Wait()

	require.Equal(t, false, err != nil)
}

func TestParse(t *testing.T) {

	urls := []string{mockXml}

	p, _ := setupParser(len(urls))
	err := p.parseAsync(urls)
	require.True(t, err == nil)

	finalItems := p.formatChannelData()
	require.Equal(t, 10, len(finalItems))

	// avoids trying to close a closed channel
	time.Sleep(100)
	badUrls := []string{
		"failMe/",
	}
	p, _ = setupParser(len(badUrls))
	err = p.parseAsync(badUrls)
	require.True(t, err != nil)
}

func Test_parseTime(t *testing.T) {
	_, err := parseTime("fail me")
	require.True(t, err != nil)
}

func Test_parseXml(t *testing.T) {
	xml, err := parseXml([]byte{0}, "")
	if err != nil {
		return
	}
	require.True(t, err != nil)
	require.True(t, xml == nil)
}
