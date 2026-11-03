package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const baseURL = "https://api.indec.gob.ar/v1/"

type IndecClient struct {
	token string
	http  *http.Client
}

func NewClient() *IndecClient {
	return &IndecClient{
		token: os.Getenv("INDEC_API_TOKEN"),
		http:  &http.Client{},
	}
}

func (c *IndecClient) doGet(path string, query url.Values) ([]byte, error) {
	u, _ := url.Parse(baseURL + path)
	if query != nil {
		u.RawQuery = query.Encode()
	}
	req, _ := http.NewRequest("GET", u.String(), nil)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("indec api error: %s - %s", resp.Status, string(b))
	}
	return b, nil
}

func main() {
	s := server.NewMCPServer(
		"indec-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	c := NewClient()

	// List datasets
	s.AddTool(mcp.NewTool("list_datasets",
		mcp.WithDescription("List available INDEC API datasets from OpenAPI spec"),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// We have cached paths from earlier
		datasets := []string{
			"ipi_manufacturero", "ipi_minero", "estimaciones_poblacion_departamento_2022_2035",
			"proyecciones_poblacion_jurisdiccion_2022_2040", "proyecciones_poblacion_nacional_2022_2040",
			"valorizacion_canastas_cba_cbt", "ipc_nacional", "redeterminacion_precios_contratos_obras_publicas",
			"sipm_precios_mayoristas", "icc_precios_construccion", "icc_indice_general", "icc_insumos",
			"icc_materiales", "icc_mano_de_obra", "icc_vivienda", "icc_edificios", "icc_obras_publicas",
			"icc_maquinaria", "icc_hormigon_elaborado",
		}
		b, _ := json.MarshalIndent(datasets, "", "  ")
		return mcp.NewToolResultText(string(b)), nil
	})

	// Query dataset
	s.AddTool(mcp.NewTool("query_dataset",
		mcp.WithDescription("Query an INDEC dataset with filters, pagination, ordering"),
		mcp.WithString("dataset", mcp.Required(), mcp.Description("Dataset name/path")),
		mcp.WithString("select", mcp.Description("Comma-separated columns to select")),
		mcp.WithString("order", mcp.Description("Order by (e.g. periodo.desc)")),
		mcp.WithNumber("limit", mcp.Description("Max rows (default 1000)")),
		mcp.WithNumber("offset", mcp.Description("Offset (default 0)")),
		mcp.WithString("filters", mcp.Description("JSON object of filters like {\"anio\": \"gt.2023\"}")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		dataset := args["dataset"].(string)
		q := url.Values{}
		if v, ok := args["select"].(string); ok && v != "" {
			q.Set("select", v)
		}
		if v, ok := args["order"].(string); ok && v != "" {
			q.Set("order", v)
		}
		limit := 1000
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		q.Set("limit", strconv.Itoa(limit))
		offset := 0
		if v, ok := args["offset"].(float64); ok {
			offset = int(v)
		}
		q.Set("offset", strconv.Itoa(offset))
		if v, ok := args["filters"].(string); ok && v != "" {
			var f map[string]string
			if json.Unmarshal([]byte(v), &f) == nil {
				for k, val := range f {
					q.Set(k, val)
				}
			}
		}
		b, err := c.doGet(dataset, q)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		// Pretty print if JSON array
		var raw json.RawMessage
		if json.Unmarshal(b, &raw) == nil {
			var out interface{}
			if json.Unmarshal(raw, &out) == nil {
				pretty, _ := json.MarshalIndent(out, "", "  ")
				return mcp.NewToolResultText(string(pretty)), nil
			}
		}
		return mcp.NewToolResultText(string(b)), nil
	})

	// Quick IPC
	s.AddTool(mcp.NewTool("get_ipc_nacional",
		mcp.WithDescription("Get national CPI (IPC nacional) data"),
		mcp.WithNumber("limit", mcp.Description("Max rows")),
		mcp.WithString("order", mcp.Description("Order (default periodo.desc)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		q := url.Values{}
		limit := 12
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		q.Set("limit", strconv.Itoa(limit))
		order := "periodo.desc"
		if v, ok := args["order"].(string); ok && v != "" {
			order = v
		}
		q.Set("order", order)
		b, err := c.doGet("ipc_nacional", q)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var out interface{}
		if json.Unmarshal(b, &out) == nil {
			pretty, _ := json.MarshalIndent(out, "", "  ")
			return mcp.NewToolResultText(string(pretty)), nil
		}
		return mcp.NewToolResultText(string(b)), nil
	})

	// Quick EMAE
	s.AddTool(mcp.NewTool("get_emae_general",
		mcp.WithDescription("Get EMAE nivel general"),
		mcp.WithNumber("limit", mcp.Description("Max rows")),
		mcp.WithString("order", mcp.Description("Order")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		q := url.Values{}
		limit := 12
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		q.Set("limit", strconv.Itoa(limit))
		order := "anio.desc,mes.desc"
		if v, ok := args["order"].(string); ok && v != "" {
			order = v
		}
		q.Set("order", order)
		b, err := c.doGet("estimador-mensual-de-actividad-economica-nivel-general", q)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var out interface{}
		if json.Unmarshal(b, &out) == nil {
			pretty, _ := json.MarshalIndent(out, "", "  ")
			return mcp.NewToolResultText(string(pretty)), nil
		}
		return mcp.NewToolResultText(string(b)), nil
	})

	// Quick EPH unemployment
	s.AddTool(mcp.NewTool("get_eph_unemployment",
		mcp.WithDescription("Get EPH mercado de trabajo principales tasas"),
		mcp.WithNumber("limit", mcp.Description("Max rows")),
		mcp.WithString("area", mcp.Description("Area geografica (optional)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		q := url.Values{}
		limit := 12
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		q.Set("limit", strconv.Itoa(limit))
		q.Set("order", "periodo.desc")
		if v, ok := args["area"].(string); ok && v != "" {
			q.Set("area_geografica_nom", "eq."+v)
		}
		b, err := c.doGet("mercado_trabajo_eph", q)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var out interface{}
		if json.Unmarshal(b, &out) == nil {
			pretty, _ := json.MarshalIndent(out, "", "  ")
			return mcp.NewToolResultText(string(pretty)), nil
		}
		return mcp.NewToolResultText(string(b)), nil
	})

	// EPH poverty
	s.AddTool(mcp.NewTool("get_eph_poverty",
		mcp.WithDescription("Get EPH pobreza e indigencia"),
		mcp.WithNumber("limit", mcp.Description("Max rows")),
		mcp.WithString("area", mcp.Description("Area geografica (optional)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		q := url.Values{}
		limit := 12
		if v, ok := args["limit"].(float64); ok {
			limit = int(v)
		}
		q.Set("limit", strconv.Itoa(limit))
		q.Set("order", "periodo.desc")
		if v, ok := args["area"].(string); ok && v != "" {
			q.Set("area_geografica_nom", "eq."+v)
		}
		b, err := c.doGet("pobreza_eph", q)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var out interface{}
		if json.Unmarshal(b, &out) == nil {
			pretty, _ := json.MarshalIndent(out, "", "  ")
			return mcp.NewToolResultText(string(pretty)), nil
		}
		return mcp.NewToolResultText(string(b)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
