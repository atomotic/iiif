package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/atomotic/iiif/internal/imageinfo"
	v3 "github.com/atomotic/iiif/v3"
)

func main() {
	id := flag.String("id", "", "Unique identifier for the manifest (e.g. book1)")
	base := flag.String("base", "", "Base URL where the manifest will be served (e.g. https://example.org/iiif)")
	title := flag.String("title", "", "Human-readable title of the manifest")
	source := flag.String("source", "", "Path to a directory of images or a PDF file to convert")
	destination := flag.String("destination", "", "Output directory; a subdirectory named <id> will be created inside it, containing the images and manifest.json")
	tiles := flag.Bool("tiles", false, "Generate IIIF image tiles for each image using vips dzsave (requires vips)")
	resolution := flag.Int("resolution", 150, "Resolution (DPI) used when converting PDF pages to images via mutool")
	flag.Parse()

	if *id == "" || *base == "" || *title == "" || *source == "" || *destination == "" {
		fmt.Fprintln(os.Stderr, "Usage: mkiiif -id <id> -base <url> -title <title> -source <dir|pdf> -destination <dir> [-tiles]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sanitizedID := sanitizeID(*id)
	if sanitizedID == "" {
		fmt.Fprintln(os.Stderr, "Error: -id is empty after sanitization")
		os.Exit(1)
	}
	if sanitizedID != *id {
		fmt.Fprintf(os.Stderr, "Warning: -id sanitized to %q\n", sanitizedID)
	}

	// destination/id is where images and manifest.json are written
	outDir := filepath.Join(*destination, sanitizedID)
	destAbs, err := filepath.Abs(*destination)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving destination path: %v\n", err)
		os.Exit(1)
	}
	outDirAbs, err := filepath.Abs(outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving output path: %v\n", err)
		os.Exit(1)
	}
	if !strings.HasPrefix(outDirAbs, destAbs+string(os.PathSeparator)) {
		fmt.Fprintln(os.Stderr, "Error: -id resolves outside the destination directory")
		os.Exit(1)
	}

	if err := os.MkdirAll(outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// If source is a PDF file, convert pages to images using mutool
	// If source is a directory, copy images into outDir
	info, err := os.Stat(*source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing source: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() && strings.ToLower(filepath.Ext(*source)) == ".pdf" {
		outPattern := filepath.Join(outDir, "page-%03d.png")
		cmd := exec.Command("mutool", "draw", "-r", fmt.Sprintf("%d", *resolution), "-F", "png", "-o", outPattern, *source)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running mutool: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := copyImages(*source, outDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error copying images: %v\n", err)
			os.Exit(1)
		}
	}

	manifest, err := v3.NewManifest(sanitizedID, *base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating manifest: %v\n", err)
		os.Exit(1)
	}
	manifest.Label = v3.Label{"none": {*title}}

	images, err := imageinfo.ScanDirectory(outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning images directory: %v\n", err)
		os.Exit(1)
	}

	if len(images) == 0 {
		fmt.Fprintln(os.Stderr, "No images found")
		os.Exit(1)
	}

	for i, img := range images {
		name := strings.TrimSuffix(filepath.Base(img.Path), filepath.Ext(img.Path))
		canvasID := fmt.Sprintf("p%d", i+1)

		if *tiles {
			serviceURL := fmt.Sprintf("%s/%s/%s", *base, sanitizedID, name)
			tileDir := filepath.Join(outDir, name)
			// vips appends the output directory basename to --id automatically,
			// so pass only the parent URL (base/id); vips will produce base/id/name in info.json
			vipsBaseURL := fmt.Sprintf("%s/%s", *base, sanitizedID)

			width, height, err := tileImage(img.Path, tileDir, vipsBaseURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error tiling %s: %v\n", img.Path, err)
				os.Exit(1)
			}

			bodyURL := serviceURL + "/full/max/0/default.jpg"
			if err := manifest.NewItem(canvasID, name, bodyURL, []int{width, height}, serviceURL); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding canvas for %s: %v\n", img.Path, err)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "Tiled %s -> %s\n", filepath.Base(img.Path), tileDir)

			if err := os.Remove(img.Path); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not remove original image %s: %v\n", img.Path, err)
			}

			// remove vips-properties.xml artifact left by vips dzsave
			_ = os.Remove(filepath.Join(outDir, "vips-properties.xml"))
		} else {
			bodyURL := fmt.Sprintf("%s/%s/%s", *base, sanitizedID, filepath.Base(img.Path))
			if err := manifest.NewItem(canvasID, name, bodyURL, []int{img.Width, img.Height}, ""); err != nil {
				fmt.Fprintf(os.Stderr, "Error adding canvas for %s: %v\n", img.Path, err)
				os.Exit(1)
			}
		}
	}

	manifestPath := filepath.Join(outDir, "manifest.json")
	if err := os.WriteFile(manifestPath, []byte(manifest.Serialize()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing manifest: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Manifest written to %s\n", manifestPath)

	indexHTML, err := renderViewer("manifest.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering viewer: %v\n", err)
		os.Exit(1)
	}
	indexPath := filepath.Join(outDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(indexHTML), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing index.html: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Viewer written to %s\n", indexPath)
}

// tileImage runs vips dzsave on imagePath, writes IIIF tiles into tileDir,
// generates full/max/0/default.jpg (not produced automatically by vips),
// and returns the image dimensions read from the generated info.json.
func tileImage(imagePath, tileDir, serviceURL string) (int, int, error) {
	cmd := exec.Command("vips", "dzsave", imagePath, tileDir, "--id", serviceURL, "--layout", "iiif3")
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return 0, 0, err
	}

	// vips dzsave does not generate full/max/0/default.jpg; create it from the original
	fullDir := filepath.Join(tileDir, "full", "max", "0")
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return 0, 0, fmt.Errorf("creating full image dir: %w", err)
	}
	fullCmd := exec.Command("vips", "copy", imagePath, filepath.Join(fullDir, "default.jpg"))
	fullCmd.Stderr = os.Stderr
	if err := fullCmd.Run(); err != nil {
		return 0, 0, fmt.Errorf("generating full image: %w", err)
	}

	infoPath := filepath.Join(tileDir, "info.json")
	data, err := os.ReadFile(infoPath)
	if err != nil {
		return 0, 0, fmt.Errorf("reading info.json: %w", err)
	}

	var info struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return 0, 0, fmt.Errorf("parsing info.json: %w", err)
	}
	if info.Width == 0 || info.Height == 0 {
		return 0, 0, fmt.Errorf("info.json has zero dimensions")
	}

	return info.Width, info.Height, nil
}

// sanitizeID replaces any character that is not a letter, digit, hyphen, or dot with _.
// This prevents path traversal and ensures the id is safe to use as a directory name.
var unsafeIDChars = regexp.MustCompile(`[^a-zA-Z0-9\-.]`)

func sanitizeID(id string) string {
	return strings.Trim(unsafeIDChars.ReplaceAllString(id, "_"), "_.")
}

func copyImages(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue
		}
		src := filepath.Join(srcDir, entry.Name())
		dst := filepath.Join(dstDir, entry.Name())
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			return err
		}
	}
	return nil
}
