# P2P_INTERVIEW
This project is backend side of p2p_interview project.


### Create docs

```swag init -g router.go -o docs/ -d internal/delivery/http,internal/delivery/http/handlers,internal/delivery/http/middleware,internal/domain/models```

### Create diagram

```goplantuml -recursive -show-aggregations ./ > diagram.puml```

