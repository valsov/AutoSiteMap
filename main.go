package main

import (
	"fmt"

	"github.com/valsov/autositemap/scraper"
)

func main() {
	s := scraper.NewSiteScraper("http://o11y.eu")
	pages := s.GetPages("/")
	fmt.Print(len(pages))
}
