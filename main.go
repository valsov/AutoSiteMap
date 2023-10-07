package main

import (
	"context"
	"flag"
	"time"

	"github.com/valsov/autositemap/scraper"
	"github.com/valsov/autositemap/visualizer"
)

var (
	domainUrl      string
	startPath      string
	timeout        int
	viewExportPath string
)

func init() {
	// Parameters parsing
	flag.StringVar(&domainUrl, "domainUrl", "", "Base url of the target website")
	flag.StringVar(&startPath, "startPath", "/", "Start path of the scraping")
	flag.IntVar(&timeout, "timeout", 120, "Timeout, in seconds, of the operations, defaults to 120 seconds (2 minutes)")
	flag.StringVar(&viewExportPath, "exportPath", "sitemap.html", "Target path to write the site map visualizer to")
}

func main() {
	flag.Parse()

	timeoutDuration := time.Duration(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration*time.Second)
	defer cancel()

	// Get site map
	s := scraper.NewSiteScraper(domainUrl)
	pages := s.GetPages(ctx, startPath)

	// Visualize
	visualizer.GenerateVisualizer(pages, viewExportPath)
}
