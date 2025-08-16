# URLX - Go URL Shortener

A simple, production-ready URL shortener written in Go, with:
- Gin web framework
- Modular storage (Postgres, in-memory)
- Zap structured logging (debug level)
- Auto-creation of database schema
- CI/CD ready for Render.com
- Middleware.io and New Relic logging ready (optional, see code)

## Live Demo

**Hosted at:** [https://urlx-iv45.onrender.com](https://urlx-iv45.onrender.com)

## Features
- Shorten URLs via API or web UI
- Redirect short URLs
- Delete short URLs
- Health check endpoint (`/healthz`)
- Modern web UI at `/`
- Debug-level logging for all operations
- Auto-creates Postgres schema on startup

## API Usage

### Shorten a URL
```
POST /shorten
Content-Type: application/json

{"url": "https://example.com"}
```
Response:
```
{"short": "abc123"}
```

### Redirect
```
GET /abc123
```
Redirects to the original URL.

### Delete
```
DELETE /delete/abc123
```

### Health Check
```
GET /healthz
```

## Local Development

1. **Clone the repo**
2. Set up Postgres and set `DATABASE_URL` env var
3. Run:
   ```sh
   go run main.go
   ```
4. Visit [http://localhost:8080](http://localhost:8080)

## Deployment

- Deploys automatically to Render.com using `render.yaml` (provisions Postgres, sets env vars, etc.)
- On first run, the required table is created automatically.

## Logging

- Uses Zap for all logs (debug/info/error)
- All storage and API operations are logged

## Note on Logging Performance

- Logging may be slower than expected because the Grafana Loki logging endpoint is in India, while the Render.com service is hosted in the USA.
- Each request blocks until the log is sent to Grafana Cloud (no batching or async logging is used).
- For production, consider batching logs or using a background worker to avoid blocking request handling.

## License
MIT
