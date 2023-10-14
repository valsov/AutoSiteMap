# Website map visualizer
Go utility to generate a map of a given website, following links in its pages.

## How to use
Available options:
- **-domainUrl**: domain to process
- **-startPath**: start path for the scraping
- **-timeout**: maximum time, in seconds, allowed to process the domain
- **-viewExportPath**: path to write the vis.js-based visualization file
- **-includeExternalLinks**: should the scraping include external links in the output (which are in all cases not followed)
- **-maxVisitsCount**: maximum number of links to follow during the process

```shell
websitemapper -domainUrl https://sample-domain.com
```

## Visualization using vis.js

Internal pages are displayed in blue, external links are yellow. This is configurable in `view.tmpl`.

Example of scraping visualization using the following parameters (maxVisitsCount was set to 1 to avoid the graph being too big):
```shell
websitemapper -domainUrl https://docker.com -maxVisitsCount 1 -includeExternalLinks true
```
![Screenshot](/screenshot.png)
