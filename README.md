# habor-cli
A cli for Harbor


### Generate mocks
```
mockgen -destination=harbor/mocks/mock_api.go -package=mocks github.com/disaster37/harbor-cli/harbor/api API,ArtifactAPI
```