package extractor

import (
	"strings"
	"testing"
)

func TestExtractFromScriptContent(t *testing.T) {
	// Sample script content with a fake URL
	scriptContent := `
	/* <![CDATA[ */
	var welcomebar_frontjs = {"ajaxurl":"http:\/\/fakeurl.invalid\/wp-admin\/admin-ajax.php","days":"Days","hours":"Hours","minutes":"Minutes","seconds":"Seconds","ajax_nonce":"c35d389da5"};
	/* ]]> */
	`

	expected := "http://fakeurl.invalid/wp-admin/admin-ajax.php"
	assets, err := extractFromScriptContent(scriptContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset, got %d", len(assets))
	}

	if assets[0] != expected {
		t.Errorf("expected asset %q, got %q", expected, assets[0])
	}
}

func TestExtractJSImportStrings(t *testing.T) {
	scriptContent := `
		import { slice } from "lodash";
		import React from 'react';
		import('./a.js');
		import("b.js").then(() => console.log('loaded'));
	`

	imports, err := extractJSImportStrings(scriptContent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(imports) != 4 {
		t.Fatalf("expected 4 imports, got %d", len(imports))
	}

	expected := []string{"lodash", "react", "./a.js", "b.js"}
	for idx, want := range expected {
		if imports[idx] != want {
			t.Fatalf("expected import[%d]=%q, got %q", idx, want, imports[idx])
		}
	}
}

func TestHTMLAssetsScriptJSImport(t *testing.T) {
	html := `<html><body><script type="module">import "./app.js"; import("./chunks/a.js");</script></body></html>`
	item := setupItem(html)

	assets, err := HTMLAssets(item)
	if err != nil {
		t.Fatalf("HTMLAssets error = %v", err)
	}

	got := make([]string, 0, len(assets))
	for _, asset := range assets {
		got = append(got, asset.Raw)
	}

	joined := strings.Join(got, "\n")
	if !strings.Contains(joined, "http://ex.com/app.js") {
		t.Fatalf("expected static import to be extracted, assets=%v", got)
	}
	if !strings.Contains(joined, "http://ex.com/chunks/a.js") {
		t.Fatalf("expected dynamic import to be extracted, assets=%v", got)
	}
}
