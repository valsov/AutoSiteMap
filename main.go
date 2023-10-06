package main

import (
	"context"
	"fmt"
	"time"

	"github.com/valsov/autositemap/scraper"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	s := scraper.NewSiteScraper("http://localhost")
	pages := s.GetPages(ctx, "/")
	fmt.Print(len(pages))
}
