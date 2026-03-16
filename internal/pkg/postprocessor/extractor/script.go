package extractor

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/internetarchive/Zeno/pkg/models"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

func extractFromScriptContent(content string) (assets []string, err error) {
	jsonContent := strings.SplitAfterN(content, "=", 2)

	if len(jsonContent) > 1 {
		var (
			openSeagullCount   int
			closedSeagullCount int
			payloadEndPosition int
		)

		// figure out the end of the payload
		for pos, char := range jsonContent[1] {
			if char == '{' {
				openSeagullCount++
			} else if char == '}' {
				closedSeagullCount++
			} else {
				continue
			}

			if openSeagullCount > 0 {
				if openSeagullCount == closedSeagullCount {
					payloadEndPosition = pos
					break
				}
			}
		}

		if len(jsonContent[1]) > payloadEndPosition {
			URLsFromJSON, _, err := GetURLsFromJSON(json.NewDecoder(strings.NewReader(jsonContent[1][:payloadEndPosition+1])))
			if err != nil {
				return nil, err
			} else {
				assets = append(assets, URLsFromJSON...)
			}
		}
	}

	return assets, nil
}

func IsJavaScript(URL *models.URL) bool {
	mimeType := URL.GetMIMEType()
	if mimeType == nil {
		return false
	}

	value := strings.ToLower(mimeType.String())
	return strings.Contains(value, "javascript") || strings.Contains(value, "ecmascript")
}

func ExtractFromURLScript(URL *models.URL) (jsImportLinks []*models.URL, otherJSLinks []*models.URL, otherCSSLinks []*models.URL, otherStrings []string, err error) {
	defer URL.RewindBody()
	bodyBytes, err := io.ReadAll(URL.GetBody())
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rawLinks, _otherStrings, err := extractJSImportStrings(string(bodyBytes))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	jsImportLinks = make([]*models.URL, 0, len(rawLinks))
	for _, rawLink := range rawLinks {
		resolvedURL, resolveErr := resolveURL(rawLink, URL)
		if resolveErr != nil || resolvedURL == "" {
			jsImportLinks = append(jsImportLinks, &models.URL{Raw: rawLink})
			continue
		}

		jsImportLinks = append(jsImportLinks, &models.URL{Raw: resolvedURL})
	}
	otherJSLinks = make([]*models.URL, 0, len(_otherStrings))
	for _, str := range _otherStrings {
		if strings.HasSuffix(str, ".css") {
			otherCSSLinks = append(otherCSSLinks, &models.URL{Raw: str})
		} else if strings.HasSuffix(str, ".js") {
			otherJSLinks = append(otherJSLinks, &models.URL{Raw: str})
		} else {
			otherStrings = append(otherStrings, str)
		}
	}

	return jsImportLinks, otherJSLinks, otherCSSLinks, otherStrings, nil
}

func extractJSImportStrings(content string) (imports []string, otherStrings []string, err error) {
	input := parse.NewInput(strings.NewReader(content))
	ast, err := js.Parse(input, js.Options{})
	if err != nil {
		return nil, nil, err
	}

	imports = make([]string, 0, 4)
	otherStrings = make([]string, 0, 4)
	js.Walk(&jsVisitor{imports: &imports, otherStrings: &otherStrings}, ast)

	return imports, otherStrings, nil
}

type jsVisitor struct {
	imports      *[]string
	otherStrings *[]string
}

func (v *jsVisitor) Enter(n js.INode) js.IVisitor {
	if call, ok := n.(*js.CallExpr); ok {
		if lit, ok := call.X.(*js.LiteralExpr); ok && lit.TokenType == js.ImportToken {
			if len(call.Args.List) > 0 {
				if strLit, ok := call.Args.List[0].Value.(*js.LiteralExpr); ok && strLit.TokenType == js.StringToken {
					*v.imports = append(*v.imports, strings.Trim(string(strLit.Data), `"'`))
				}
			}
		}
	}

	if importStmt, ok := n.(*js.ImportStmt); ok {
		*v.imports = append(*v.imports, strings.Trim(string(importStmt.Module), `"'`))
	}

	if lit, ok := n.(*js.LiteralExpr); ok && lit.TokenType == js.StringToken {
		*v.otherStrings = append(*v.otherStrings, strings.Trim(string(lit.Data), `"'`))
	}

	return v
}

func (v *jsVisitor) Exit(js.INode) {}
