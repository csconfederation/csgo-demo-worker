# csgo-demo-worker

## How to Use

### Set Environment Variables

- `LOG_LEVEL` - logrus log level
- `DEMO_STATS_USER` - username for basic auth (optional)
- `DEMO_STATS_PASSWORD` - password for basic auth (optional)

### Endpoints

| Path               | Method | Body               | Parameters                                                        |
|--------------------|--------|--------------------|-------------------------------------------------------------------|
| `api/parse`        | POST   | Binary `.dem` file | n/a                                                               |
| `api/parse-remote` | GET    | n/a                | `url` - remote url, `auth` - Full Authorization header (optional) |

### Docker
```bash
docker run \
-p 8080:8080 \ 
-e PORT=8080 \
ghcr.io/csconfederation/csgo-demo-worker:latest
```


## Libraries Used

- [gin-gonic](https://github.com/gin-gonic/) - web server
- [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang) - base library for demo parsing
