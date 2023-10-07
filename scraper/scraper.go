package scraper

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"
)

type Page struct {
	Id                    int
	Url                   string
	OutgoingLinks         []int
	IsInternalUrl, Failed bool
}

type SiteScraper struct {
	baseUrl *url.URL
	db      map[string]*Page
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewSiteScraper(domainUrl string) *SiteScraper {
	parsedUrl, err := url.Parse(domainUrl)
	if err != nil {
		panic(err)
	}
	return &SiteScraper{
		baseUrl: parsedUrl,
		db:      map[string]*Page{},
		mutex:   sync.Mutex{},
		wg:      sync.WaitGroup{},
	}
}

func (s *SiteScraper) GetPages(ctx context.Context, path string) []*Page {
	initialPage := &Page{Id: 1, Url: path, OutgoingLinks: []int{}, IsInternalUrl: true}
	s.db[path] = initialPage

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.visit(ctx, initialPage)
	}()
	s.wg.Wait()

	pages := make([]*Page, 0, len(s.db))
	for _, page := range s.db {
		pages = append(pages, page)
	}
	return pages
}

func (s *SiteScraper) visit(ctx context.Context, page *Page) {
	s.wg.Add(1)
	defer s.wg.Done()

	requestUrl := s.baseUrl.JoinPath(page.Url)
	links, err := getLinksFromPage(ctx, requestUrl.String())
	if err != nil {
		page.Failed = true
		return
	}

	for _, link := range links {
		link, internal := s.formatInternalUrl(link) // Format link
		s.mutex.Lock()
		if linkedPage, found := s.db[link]; found {
			page.OutgoingLinks = append(page.OutgoingLinks, linkedPage.Id)
		} else {
			linkedPage = &Page{Id: len(s.db) + 1, Url: link, OutgoingLinks: []int{}, IsInternalUrl: internal}
			page.OutgoingLinks = append(page.OutgoingLinks, linkedPage.Id)
			s.db[link] = linkedPage // Add to cache
			if linkedPage.IsInternalUrl {
				// Only visit internal resources
				go s.visit(ctx, linkedPage)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *SiteScraper) formatInternalUrl(inputUrl string) (string, bool) {
	parsed, err := url.Parse(inputUrl)
	if err != nil {
		return inputUrl, false
	}

	if parsed.Host == s.baseUrl.Host {
		return parsed.Path, true // Ommit query string and host
	}
	return inputUrl, false
}

func getLinksFromPage(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
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
			if len(key) == 4 && string(key) == "href" {
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
