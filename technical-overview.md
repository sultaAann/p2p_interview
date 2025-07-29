### Create docs

```swag init -g router.go -o docs/ -d internal/delivery/http,internal/delivery/http/handlers,internal/delivery/http/middleware,internal/domain/models```

### Create diagram

```goplantuml -recursive -show-aggregations ./ > diagram.puml```

