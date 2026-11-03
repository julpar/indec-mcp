# INDEC MCP Server

Servidor MCP para acceder a los datos estadísticos oficiales del [INDEC Argentina](https://www.indec.gob.ar/) a través de su [API pública](https://api.indec.gob.ar/).

Este servidor permite que asistentes de IA como Claude, Cursor, Cline u otros clientes compatibles con MCP consulten directamente series históricas como IPC, EMAE, EPH, comercio exterior y mucho más.

## Características

- **Consulta genérica**: Accede a cualquier dataset disponible en la API del INDEC con filtros, ordenamiento y paginación.
- **Acceso rápido**: Herramientas listas para usar para las series más consultadas (IPC, EMAE, mercado laboral y pobreza).
- **Filtrado tipo PostgREST**: Soporta operadores como `eq.`, `gt.`, `gte.`, `lt.`, `lte.`, `ilike.*...*`, entre otros.
- **Paginación**: Control total sobre `limit` y `offset`.
- **Autenticación con token**: Compatible con el esquema Bearer requerido por la API.

## Instalación

### Requisitos

- [Go](https://go.dev/) 1.21 o superior (solo necesario si querés compilar desde código)
- Un token de acceso a la API del INDEC (configurado como variable de entorno)

### Compilar desde código

```bash
git clone https://github.com/julpar/indec-mcp.git
cd indec-mcp
go build -o indec-mcp .
```

También podés usar el binario ya incluido en el repositorio.

## Configuración

La API requiere autenticación. Configurá tu token como variable de entorno:

```bash
export INDEC_API_TOKEN="tu_token_aqui"
```

### Configurar en Claude Desktop

Editá `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "indec": {
      "command": "/ruta/completa/a/indec-mcp",
      "env": {
        "INDEC_API_TOKEN": "tu_token_aqui"
      }
    }
  }
}
```

### Configurar en Claude Code

```bash
claude mcp add indec /ruta/completa/a/indec-mcp --env INDEC_API_TOKEN=tu_token_aqui
```

### Configurar en Cursor / Cline

Agregalo en la configuración de MCP servers con el mismo formato (comando + variable de entorno).

## Herramientas disponibles

| Herramienta | Descripción | Parámetros principales |
|-------------|-------------|------------------------|
| `list_datasets` | Lista algunos de los datasets más comunes del INDEC. | Ninguno |
| `query_dataset` | Consulta genérica a cualquier dataset de la API. | `dataset`, `select`, `order`, `limit`, `offset`, `filters` (JSON) |
| `get_ipc_nacional` | Obtiene datos del Índice de Precios al Consumidor (IPC) nacional. | `limit`, `order` |
| `get_emae_general` | Obtiene el Estimador Mensual de Actividad Económica (EMAE) a nivel general. | `limit`, `order` |
| `get_eph_unemployment` | Obtiene las principales tasas del mercado laboral (EPH). | `limit`, `area` |
| `get_eph_poverty` | Obtiene incidencia de pobreza e indigencia (EPH). | `limit`, `area` |

## Ejemplos de uso

Una vez conectado el servidor MCP, podés hacer preguntas como:

- *"Mostrame la inflación (IPC) de los últimos 12 meses"*
- *"Cómo viene la actividad económica según el EMAE?"*
- *"Cuál es la tasa de desocupación según la EPH?"*
- *"Dame la evolución de la pobreza en el Gran Buenos Aires"*
- *"Consultame los datos de comercio exterior de los últimos meses"*

## Notas técnicas

- **URL base**: `https://api.indec.gob.ar/v1/`
- **Estilo PostgREST**: Los filtros se pasan tal como lo espera la API (ej. `{"anio": "gt.2024"}`).
- **Paginación**: Por defecto usa `limit=1000`, `offset=0`. Para obtener totales podés usar el header `Prefer: count=exact` desde `query_dataset` si fuera necesario.
- **Autenticación**: Requiere token Bearer. Sin token, la API devuelve `permission denied` para la mayoría de los endpoints.
- **Transporte**: stdio (compatible con la gran mayoría de clientes MCP).

## Licencia

MIT
