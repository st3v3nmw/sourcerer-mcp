package parser_test

import (
	"path/filepath"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
	"github.com/stretchr/testify/suite"
)

type ParserBaseTestSuite struct {
	suite.Suite
	parser        *parser.Parser
	workspaceRoot string
}

func (s *ParserBaseTestSuite) SetupSuite() {
	s.workspaceRoot = filepath.Join("..", "..", "testdata")
}

func (s *ParserBaseTestSuite) getChunks(filePath string) map[string]*parser.Chunk {
	file, err := s.parser.Chunk(filePath)
	s.Require().NoError(err)
	s.Require().NotNil(file)

	chunks := make(map[string]*parser.Chunk)
	for _, chunk := range file.Chunks {
		chunks[chunk.Path] = chunk

		s.T().Logf("%s | %s", chunk.Path, chunk.Summary)
	}

	return chunks
}

func (s *GoParserTestSuite) TearDownSuite() {
	if s.parser != nil {
		s.parser.Close()
	}
}
