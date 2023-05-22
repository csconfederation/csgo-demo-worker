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

## Deploy Checklist

Most changes to this repo will be updating the parser lib version, here is a checklist for this process:

- [ ] Update [parser](https://github.com/csconfederation/demoScrape2) version in `go.mod`
- [ ] Create release with a new tag using [semantic versioning](https://semver.org/) (this kicks off deploy action)

## Libraries Used

- [gin-gonic](https://github.com/gin-gonic/) - web server
- [demoinfocs-golang](https://github.com/markus-wa/demoinfocs-golang) - base library for demo parsing
