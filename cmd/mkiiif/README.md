# mkiiif

Generates a IIIF Presentation API v3 manifest from a directory of images or a PDF file.

## Usage

```
mkiiif -id <id> -base <url> -title <title> -source <dir|pdf> -destination <dir> [options]
```

**Required flags**

| Flag | Description |
|------|-------------|
| `-id` | Identifier for the manifest; used as output subdirectory name |
| `-base` | Base URL where the output will be served (e.g. `https://example.org/iiif`) |
| `-title` | Human-readable title |
| `-source` | Path to an image directory or a PDF file |
| `-destination` | Directory where `<id>/` will be created |

**Optional flags**

| Flag | Default | Description |
|------|---------|-------------|
| `-resolution` | `150` | DPI used when rasterizing PDF pages via `mutool` |
| `-tiles` | `false` | Generate IIIF image tiles using `vips dzsave` |

## Output

```
<destination>/<id>/
├── manifest.json   # IIIF manifest
├── index.html      # Triiiceratops viewer
└── *.png / *.jpg   # images (or tile directories if -tiles is set)
```

## Dependencies

- [`mutool`](https://mupdf.com/) — required for PDF input
- [`vips`](https://www.libvips.org/) — required when using `-tiles`

## Example

```sh
# From a PDF
mkiiif -id book1 -base https://example.org/iiif -title "My Book" \
       -source document.pdf -destination /var/www/iiif

# From an image directory, with tiling
mkiiif -id book1 -base https://example.org/iiif -title "My Book" \
       -source ./scans -destination /var/www/iiif -tiles
```
