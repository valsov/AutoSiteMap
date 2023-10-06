package scraper

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Page struct {
	Url                 string
	OutgoingLinks       []string
	InternalUrl, Failed bool
}

type SiteScraper struct {
	baseUrl string
	db      map[string]*Page
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewSiteScraper(domainUrl string) *SiteScraper {
	return &SiteScraper{
		baseUrl: domainUrl,
		db:      map[string]*Page{},
		mutex:   sync.Mutex{},
		wg:      sync.WaitGroup{},
	}
}

func (s *SiteScraper) GetPages(path string) []*Page {
	initialPage := &Page{Url: path, OutgoingLinks: []string{}}
	s.db[s.baseUrl] = initialPage

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.visit(initialPage)
	}()
	s.wg.Wait()

	pages := make([]*Page, 0, len(s.db))
	for _, page := range s.db {
		pages = append(pages, page)
	}
	return pages
}

func (s *SiteScraper) visit(page *Page) {
	s.wg.Add(1)
	defer s.wg.Done()

	links, err := getLinksFromPage(s.baseUrl + page.Url) // todo: use net/url
	if err != nil {
		page.Failed = true
	}

	for _, link := range links {
		link, internal := s.formatInternalUrl(link) // Format link
		s.mutex.Lock()
		if linkedPage, found := s.db[link]; found {
			page.OutgoingLinks = append(page.OutgoingLinks, linkedPage.Url)
		} else {
			linkedPage = &Page{Url: link, OutgoingLinks: []string{}, InternalUrl: internal}
			page.OutgoingLinks = append(page.OutgoingLinks, linkedPage.Url)
			s.db[link] = linkedPage // Add to cache
			if linkedPage.InternalUrl {
				// Only visit internal resources
				go s.visit(linkedPage)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *SiteScraper) formatInternalUrl(url string) (string, bool) {
	// todo: use net/url
	if len(url) == 0 {
		return url, false
	}
	if url[0] == '/' {
		return url, true
	}
	if strings.HasPrefix(url, s.baseUrl) {
		return url[len(s.baseUrl):], true
	}
	return url, false
}

func getLinksFromPage(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	links := []string{}
	tokenizer := html.NewTokenizer(res.Body)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt != html.StartTagToken {
			continue
		}

		tName, hasAttr := tokenizer.TagName()
		if !hasAttr || len(tName) != 1 || tName[0] != 'a' {
			continue
		}

		for {
			key, val, hasMoreAttr := tokenizer.TagAttr()
			if len(key) == 4 && key[0] == 'h' && key[1] == 'r' && key[2] == 'e' && key[3] == 'f' {
				links = append(links, string(val))
				break
			}
			if !hasMoreAttr {
				break
			}
		}
	}

	return links, nil
}
