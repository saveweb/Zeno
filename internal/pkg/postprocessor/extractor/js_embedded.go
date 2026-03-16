package extractor

import (
	"github.com/internetarchive/Zeno/internal/pkg/config"
	"github.com/internetarchive/Zeno/internal/pkg/log"
	"github.com/internetarchive/Zeno/pkg/models"
)

var jsLogger = log.NewFieldedLogger(&log.Fields{
	"component": "postprocessor.extractor.js",
})

// IsEmbeddedJavaScript checks if the item is an embedded JavaScript module.
// An embedded JavaScript module is a JavaScript item linked from an HTML item,
// potentially through chained JavaScript imports.
func IsEmbeddedJavaScript(item *models.Item) bool {
	ok, _ := isEmbeddedJavaScriptWithJump(item, 0)
	return ok
}

// Returns the number of JS import jumps to reach the HTML item.
//
// for example:
//
//	JS -> HTML: 0 import jump
//	JS -> JS -> HTML: 1 import jump
//	JS -> JS -> JS -> HTML: 2 import jump
func GetEmbeddedJavaScriptJump(item *models.Item) int {
	ok, importJump := isEmbeddedJavaScriptWithJump(item, 0)
	if !ok {
		jsLogger.Warn("item is not embedded JavaScript, returning 0 import jump", "func", "GetEmbeddedJavaScriptJump", "item_id", item.GetShortID())
		return 0
	}
	return importJump
}

func isEmbeddedJavaScriptWithJump(item *models.Item, importJump int) (bool, int) {
	base := skipRedirectedItem(item)
	if base == nil {
		return false, 0
	}

	if IsJavaScript(base.GetURL()) {
		parent := skipRedirectedItem(base.GetParent())
		if parent == nil {
			return false, 0
		}

		if IsJavaScript(parent.GetURL()) {
			return isEmbeddedJavaScriptWithJump(parent, importJump+1)
		} else if IsHTML(parent.GetURL()) {
			return true, importJump
		}

		return false, 0
	}

	return false, 0
}

// Add JSImport links to the item as children if import jump is less than --max-js-jump.
// We reuse --max-js-jump as a safety cap for chained module imports.
func AddJSImportLinksToItemChild(item *models.Item, jsImportLinks []*models.URL) {
	if len(jsImportLinks) == 0 {
		return
	}

	if GetEmbeddedJavaScriptJump(item) >= config.Get().MaxJSJump {
		jsLogger.Warn("item is embedded JavaScript with import jump more than --max-js-jump, discarding js_import_links", "item_id", item.GetShortID(), "max_jump", config.Get().MaxJSJump)
		return
	}

	for _, link := range jsImportLinks {
		newURL := &models.URL{
			Raw:  link.Raw,
			Hops: item.GetURL().GetHops(),
		}

		newChild := models.NewItem(newURL, "")
		err := item.AddChild(newChild, models.ItemGotChildren)
		if err != nil {
			panic(err)
		}
	}
}
