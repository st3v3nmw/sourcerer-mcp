## semantic_search

{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "semantic_search", "arguments": {"query":"How is code parsed in this codebase?"}}}

{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "semantic_search", "arguments": {"query":"Where is the database code that handles SQL queries?"}}}

{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "semantic_search", "arguments": {"query":"How does authentication work in this codebase?"}}}

{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "semantic_search", "arguments": {"query":"What functions handle parsing Go code?"}}}

## get_tocs

{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "get_tocs", "arguments": {"paths": ["internal/parser/go.go"]}}}

## get_source_code

{"jsonrpc": "2.0", "id": 1, "method": "tools/call", "params": {"name": "get_source_code", "arguments": {"paths": ["internal/parser/parser.go::ParserBase::getTextWithQuery"]}}}
