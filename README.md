# iiif

> **Work in progress.**

Go library for building [IIIF Presentation API v3](https://iiif.io/api/presentation/3.0/) manifests, and a CLI tool ([mkiiif](cmd/mkiiif/README.md)) for generating them from image directories or PDF files.

## Install

```sh
go get github.com/atomotic/iiif
```

## Usage

```go
package main

import (
	"fmt"
	v3 "github.com/atomotic/iiif/v3"
)

func main() {
	manifest, _ := v3.NewManifest("book1", "https://example.org/iiif")
	manifest.Label = v3.Label{"en": {"Book 1"}}

	manifest.NewItem(
		"p1",                    // canvas id
		"Page 1",                // label
		"https://example.org/iiif/book1/page1/full/max/0/default.jpg", // image URL
		[]int{1500, 2000},       // width, height
		"",                      // image service URL (optional)
	)

	fmt.Println(manifest.Serialize())
}
```

When an [IIIF Image API](https://iiif.io/api/image/3.0/) service URL is provided as the last argument to `NewItem`, dimensions are fetched automatically from its `info.json` and a service reference is added to the canvas body:

```go
manifest.NewItem("p1", "Page 1", bodyURL, nil, "https://example.org/iiif/book1/page1")
```
