# indec-mcp

MCP server for [INDEC Argentina](https://api.indec.gob.ar/) official statistics API (PostgREST). Provides tools to query IPC, EMAE, EPH (unemployment/poverty), and any dataset from the INDEC API.

## Features

- `list_datasets` - List common INDEC datasets
- `query_dataset` - Query any dataset with filters, select, order, pagination
- `get_ipc_nacional` - Quick IPC (CPI) access
- `get_emae_general` - Quick EMAE (economic activity estimator) access
- `get_eph_unemployment` - Quick EPH unemployment rates
- `get_eph_poverty` - Quick EPH poverty/indigence rates

## Requirements

- Go 1.21+
- INDEC API token (set via `INDEC_API_TOKEN` env var). The API requires Bearer auth as per OpenAPI spec.

## Build

```bash
cd indec-mcp
go build .
```

## Usage (stdio MCP)

Add to your MCP config (e.g. Claude Desktop, Claude Code, Cursor):

```json
{
  "mcpServers": {
    "indec": {
      "command": "/path/to/indec-mcp/indec-mcp",
      "env": {
        "INDEC_API_TOKEN": "your-token-here"
      }
    }
  }
}
```

## Notes

- Base URL: `https://api.indec.gob.ar/v1/`
- PostgREST-style filtering: use operators like `eq.`, `gt.`, `lt.`, `gte.`, `lte.`, `ilike.*value*`
- Pagination uses `limit`/`offset`, with `Prefer: count=exact` for totals
- Many endpoints require auth (as observed)

## License

MIT
