package main

import (
	"context"
	"flag"
	"time"

	"github.com/valsov/websitemapper/scraper"
	"github.com/valsov/websitemapper/visualizer"
)

var (
	domainUrl            string
	startPath            string
	timeout              int
	viewExportPath       string
	includeExternalLinks bool
	maxVisitsCount       int
)

func init() {
	// Parameters parsing
	flag.StringVar(&domainUrl, "domainUrl", "", "Base url of the target website")
	flag.StringVar(&startPath, "startPath", "/", "Start path of the scraping")
	flag.IntVar(&timeout, "timeout", 120, "Timeout, in seconds, of the operations, defaults to 120 seconds (2 minutes)")
	flag.StringVar(&viewExportPath, "exportPath", "sitemap.html", "Target path to write the site map visualizer to")
	flag.BoolVar(&includeExternalLinks, "includeExternalLinks", false, "Should the scraping include external links in the output (which are in all cases not followed)")
	flag.IntVar(&maxVisitsCount, "maxVisitsCount", 0, "Maximum number of links to follow during the process")
}

func main() {
	flag.Parse()

	timeoutDuration := time.Duration(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration*time.Second)
	defer cancel()

	// Get site map
	s := scraper.NewSiteScraper(domainUrl, maxVisitsCount, includeExternalLinks)
	pages := s.GetPages(ctx, startPath)

	// Visualize
	visualizer.GenerateVisualizer(pages, viewExportPath)
}
