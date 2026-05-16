package engine

import (
	"strings"
	"testing"

	"gridea-pro/backend/internal/domain"
)

func newPP(katexEnabled bool) *HtmlPostProcessor {
	return NewHtmlPostProcessor(
		&domain.SeoSetting{}, &domain.CdnSetting{}, &domain.PwaSetting{},
		"https://example.com",
		"Site", "Desc", "zh-CN",
		"", "default", "1.0.0",
		katexEnabled,
	)
}

// TestInjectKatexCSS_WithMath 含 KaTeX 公式且开关开 → 注入 katex.min.css 到 head。
func TestInjectKatexCSS_WithMath(t *testing.T) {
	in := `<html><head><title>x</title></head><body><p><span class="katex"><math></math></span></p></body></html>`
	out := newPP(true).injectKatexCSS(in)
	if !strings.Contains(out, "katex.min.css") {
		t.Fatalf("含公式的页面没注入 KaTeX CSS:\n%s", out)
	}
}

// TestInjectKatexCSS_NoMath 不含公式 → 不注入（节省一个 stylesheet 请求）。
func TestInjectKatexCSS_NoMath(t *testing.T) {
	in := `<html><head><title>x</title></head><body><p>hello</p></body></html>`
	out := newPP(true).injectKatexCSS(in)
	if strings.Contains(out, "katex.min.css") {
		t.Fatalf("无公式的页面不应该注入 KaTeX CSS:\n%s", out)
	}
}

// TestInjectKatexCSS_Disabled 开关关 → 即使含公式也不注入（由主题自行负责）。
func TestInjectKatexCSS_Disabled(t *testing.T) {
	in := `<html><head><title>x</title></head><body><p><span class="katex"></span></p></body></html>`
	out := newPP(false).injectKatexCSS(in)
	if strings.Contains(out, "katex.min.css") {
		t.Fatalf("开关关闭时不应该注入 KaTeX CSS:\n%s", out)
	}
}

// TestInjectKatexCSS_AlreadyHasLink 主题已经引了 katex CSS → 不重复注入。
func TestInjectKatexCSS_AlreadyHasLink(t *testing.T) {
	in := `<html><head><link href="https://example.com/katex.min.css"><title>x</title></head><body><span class="katex"></span></body></html>`
	out := newPP(true).injectKatexCSS(in)
	if strings.Count(out, "katex.min.css") != 1 {
		t.Fatalf("主题已引 katex CSS 时不应再注入:\n%s", out)
	}
}
