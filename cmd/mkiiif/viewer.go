package main

import (
	"bytes"
	"text/template"
)

var viewerTemplate = template.Must(template.New("viewer").Parse(`<!doctype html>
<html>
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>Triiiceratops Viewer</title>
        <style>
            * {
                margin: 0;
                padding: 0;
                box-sizing: border-box;
            }
            html,
            body {
                width: 100%;
                height: 100%;
                overflow: hidden;
            }
            triiiceratops-viewer {
                display: block;
                width: 100%;
                height: 100%;
            }
        </style>
    </head>
    <body>
        <triiiceratops-viewer
            manifest-id="{{ .ManifestID }}"
            theme="light"
            config='{
              "renderer": {
                "maxZoomPixelRatio": 4,
                "zoomPerWheelNotch": 1.5
              },
              "showToggle": true,
              "toolbarOpen": false,
              "showCanvasNav": true,
              "showZoomControls": true,
              "viewingMode": "individuals",
              "toolbar": {
                "showSearch": false,
                "showGallery": true,
                "showAnnotations": false,
                "showFullscreen": true,
                "showInfo": true,
                "showViewingMode": true
              },
              "gallery": {
                "open": true,
                "dockPosition": "bottom",
                "size": 102
              },
              "search": {
                "open": false,
                "showCloseButton": true,
                "query": ""
              },
              "annotations": {
                "open": false
              },
              "transparentBackground": false
            }'
        >
        </triiiceratops-viewer>
        <script src="https://unpkg.com/triiiceratops@^1/dist/triiiceratops-element.iife.js"></script>
    </body>
</html>
`))

func renderViewer(manifestID string) (string, error) {
	var buf bytes.Buffer
	err := viewerTemplate.Execute(&buf, struct{ ManifestID string }{manifestID})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
