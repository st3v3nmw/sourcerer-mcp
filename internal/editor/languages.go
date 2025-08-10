package editor

import (
	"fmt"
	"path"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Language string

const (
	Go          Language = "go"
	UnknownLang Language = "unknown"
)

var extensions = map[string]Language{
	".go": Go,
}

func detectLanguage(filePath string) Language {
	language, exists := extensions[path.Ext(filePath)]
	if !exists {
		return UnknownLang
	}

	return language
}

func newParserForLanguage(language Language, workspaceRoot string) (parser.Parser, error) {
	switch language {
	case Go:
		return parser.NewGoParser(workspaceRoot)
	default:
		return nil, fmt.Errorf("language %s not supported", language)
	}
}
