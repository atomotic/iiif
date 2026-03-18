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
        <link
            rel="stylesheet"
            href="https://unpkg.com/triiiceratops/dist/triiiceratops-element.css"
        />
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
                width: 100%;
                height: 100%;
                /*display: block;*/
            }
        </style>
    </head>
    <body>
        <triiiceratops-viewer
            manifest-id="{{ .ManifestID }}"
            canvas-id=""
            config='{
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
                "draggable": false,
                "showCloseButton": false,
                "dockPosition": "bottom",
                "fixedHeight": 102
              },
              "search": {
                "open": false,
                "showCloseButton": true,
                "query": ""
              },
              "annotations": {
                "open": false,
                "visible": false
              },
              "transparentBackground": false
            }'
        >
        </triiiceratops-viewer>
        <script src="https://unpkg.com/triiiceratops/dist/triiiceratops-element.iife.js"></script>
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
