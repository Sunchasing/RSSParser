package rss_reader

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// todo: test more
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

func Test_parse(t *testing.T) {
	mockDir, _ := os.Getwd()
	mockDir = filepath.Dir(mockDir)
	mockDir = filepath.Dir(mockDir) + "\\mocks\\mockedXML.xml"
	testParseToChannel(t, []string{mockDir})
}

func testParseToChannel(t *testing.T, urls []string) {
	p, wg := setupParser(len(urls))
	readBytes, _ := p.hr.get(urls[0])

	assert.Equal(t, 117823, len(readBytes))

	assert.Equal(t, &wg, p.wg)

	p.wg.Add(1)

	err := p.parseToChannel(urls[0])
	p.wg.Wait()

	assert.Equal(t, false, err != nil)
}

// todo: better
//func TestParse(t *testing.T) {
//
//	mockDir, _ := os.Getwd()
//	mockDir = filepath.Dir(mockDir)
//	mockDir = filepath.Dir(mockDir) + "\\mocks\\mockedXML.xml"
//	urls := []string{mockDir}
//
//	var rv []RssItem
//	p, _ := setupParser(len(urls))
//	_ = p.parseAsync(urls)
//
//	rv = p.formatChannelData()
//
//	numRssItems := assert.Equal(t, len(urls), len(rv))
//	if !numRssItems{
//		t.Fatalf("expected size of the return value is %v, got %v", len(urls), len(rv))
//	}
//}
