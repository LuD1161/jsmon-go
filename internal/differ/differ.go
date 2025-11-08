package differ

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/ditashi/jsbeautifier-go/jsbeautifier"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// Differ handles generating diffs between file versions
type Differ struct{}

// New creates a new Differ
func New() *Differ {
	return &Differ{}
}

// GenerateHTMLDiff creates an HTML diff between old and new content
func (d *Differ) GenerateHTMLDiff(oldContent, newContent string) (string, error) {
	// Beautify JavaScript code
	options := jsbeautifier.DefaultOptions()
	oldBeautified, err := jsbeautifier.Beautify(&oldContent, options)
	if err != nil {
		// If beautification fails, use original content
		oldBeautified = oldContent
	}

	newBeautified, err := jsbeautifier.Beautify(&newContent, options)
	if err != nil {
		// If beautification fails, use original content
		newBeautified = newContent
	}

	// Generate diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldBeautified, newBeautified, false)

	// Create HTML output
	html := d.generateHTMLFromDiffs(diffs, oldBeautified, newBeautified)
	return html, nil
}

// generateHTMLFromDiffs creates a styled HTML page showing the diff
func (d *Differ) generateHTMLFromDiffs(diffs []diffmatchpatch.Diff, oldContent, newContent string) string {
	var buf bytes.Buffer

	// HTML header with styles
	buf.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>JSMon Diff</title>
    <style>
        body {
            font-family: 'Courier New', monospace;
            font-size: 12px;
            margin: 20px;
            background-color: #f5f5f5;
        }
        .container {
            display: flex;
            gap: 20px;
        }
        .column {
            flex: 1;
            background: white;
            padding: 10px;
            border: 1px solid #ddd;
            overflow-x: auto;
        }
        .header {
            font-weight: bold;
            margin-bottom: 10px;
            padding: 5px;
            background: #e0e0e0;
        }
        .line {
            padding: 2px 5px;
            white-space: pre-wrap;
            word-break: break-all;
        }
        .delete {
            background-color: #ffcccc;
            color: #cc0000;
        }
        .insert {
            background-color: #ccffcc;
            color: #00cc00;
        }
        .equal {
            color: #333;
        }
        pre {
            margin: 0;
            font-family: inherit;
        }
    </style>
</head>
<body>
    <h1>JSMon Diff Report</h1>
    <div class="container">
        <div class="column">
            <div class="header">Old Version</div>
            <div class="content">
`)

	// Split contents into lines for side-by-side comparison
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	// Generate old version content
	for _, line := range oldLines {
		buf.WriteString(fmt.Sprintf(`                <div class="line"><pre>%s</pre></div>
`, html.EscapeString(line)))
	}

	buf.WriteString(`            </div>
        </div>
        <div class="column">
            <div class="header">New Version</div>
            <div class="content">
`)

	// Generate new version content
	for _, line := range newLines {
		buf.WriteString(fmt.Sprintf(`                <div class="line"><pre>%s</pre></div>
`, html.EscapeString(line)))
	}

	buf.WriteString(`            </div>
        </div>
    </div>
    <div style="margin-top: 30px;">
        <h2>Changes Summary</h2>
        <div style="background: white; padding: 15px; border: 1px solid #ddd;">
`)

	// Generate changes summary with highlighting
	for _, diff := range diffs {
		text := html.EscapeString(diff.Text)
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			buf.WriteString(fmt.Sprintf(`            <span class="delete">%s</span>`, text))
		case diffmatchpatch.DiffInsert:
			buf.WriteString(fmt.Sprintf(`            <span class="insert">%s</span>`, text))
		case diffmatchpatch.DiffEqual:
			// Only show a snippet of equal text to reduce clutter
			if len(text) > 100 {
				text = text[:50] + "..." + text[len(text)-50:]
			}
			buf.WriteString(fmt.Sprintf(`            <span class="equal">%s</span>`, text))
		}
	}

	buf.WriteString(`
        </div>
    </div>
</body>
</html>`)

	return buf.String()
}
