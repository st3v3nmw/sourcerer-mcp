package analyzer

import (
	"fmt"
	"path/filepath"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Language string

const (
	Go          Language = "go"
	JavaScript  Language = "javascript"
	Markdown    Language = "markdown"
	Python      Language = "python"
	UnknownLang Language = "unknown"
)

type ParserFactory func(workspaceRoot string) (*parser.Parser, error)

type registry struct {
	extensions map[string]Language
	factories  map[Language]ParserFactory
}

func (r *registry) supportedExts() []string {
	extensions := make([]string, 0, len(r.extensions))
	for ext := range r.extensions {
		extensions = append(extensions, ext)
	}

	return extensions
}

func (r *registry) detect(filePath string) Language {
	lang, exists := r.extensions[filepath.Ext(filePath)]
	if !exists {
		return UnknownLang
	}

	return lang
}

func (r *registry) createParser(workspaceRoot string, lang Language) (*parser.Parser, error) {
	factory, exists := r.factories[lang]
	if !exists {
		return nil, fmt.Errorf("language %s not supported", lang)
	}

	return factory(workspaceRoot)
}

func (r *registry) register(lang Language, extensions []string, factory ParserFactory) {
	r.factories[lang] = factory
	for _, ext := range extensions {
		r.extensions[ext] = lang
	}
}

var languages = &registry{
	extensions: map[string]Language{},
	factories:  map[Language]ParserFactory{},
}

func init() {
	languages.register(
		Go,
		[]string{".go"},
		func(workspaceRoot string) (*parser.Parser, error) {
			return parser.NewGoParser(workspaceRoot)
		},
	)

	languages.register(
		JavaScript,
		[]string{".js", ".jsx", ".mjs"},
		func(workspaceRoot string) (*parser.Parser, error) {
			return parser.NewJavaScriptParser(workspaceRoot)
		},
	)

	languages.register(
		Markdown,
		[]string{".md", ".markdown"},
		func(workspaceRoot string) (*parser.Parser, error) {
			return parser.NewMarkdownParser(workspaceRoot)
		},
	)

	languages.register(
		Python,
		[]string{".py"},
		func(workspaceRoot string) (*parser.Parser, error) {
			return parser.NewPythonParser(workspaceRoot)
		},
	)
}
