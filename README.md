[![build](https://github.com/disaster37/harbor-cli/actions/workflows/workflow.yaml/badge.svg)](https://github.com/disaster37/harbor-cli/actions/workflows/workflow.yaml)
[![GoDoc](https://godoc.org/github.com/disaster37/harbor-cli?status.svg)](http://godoc.org/github.com/disaster37/harbor-cli)
[![codecov](https://codecov.io/gh/disaster37/harbor-cli/branch/main/graph/badge.svg)](https://codecov.io/gh/disaster37/harbor-cli/branch/main)

# habor-cli
A cli for Harbor to interfact during CI/CD jobs.
We use it to force scan and check vulnerabilities report during CI before to start CD.
We also use it to delete artifact (docker image) from Harbor when volatile environment is destroyed on PR steps.

> It's work for API v2.


## CLI

### Global options

The following parameters are available for all commands line :

- **--debug**: Enable the debug mode
- **--help**: Display help for the current command
- **--url** (required): The Harbor url with suffixe API version (https://harbor.company.com/api/v2.0)
- **--username** (required): The username to connect on Harbor API
- **--password** (required): The password to connect on Harbor API
- **--timeout**: Wait time before exit on error. Default to `60s`.
- **--self-signed-certificat**: To not check validity of Harbor certificate. Default to `false`.
- **--no-color**: To not display logs with color. Default to `false`.

You can set also this parameters on yaml file (one or all) and use the parameters `--config` with the path of your Yaml file.

```yaml
---
url: https://harbor.company.com/api/v2.0
username: admin
password: admin
self-signed-certificat: true
timeout: 180s
```

### Check vulnerabilities

It permit to get vulnerabilities from Harbor on docker image and check the severity.
It also display the report as output like trivy.
If severity is biggest that the provided, it return with exit code `1`.

We use it on Jenkins CI step just after build docker image with Kaniko.

Sample of command:

```bash
harbor-cli --url https://harbor.company.com/api/v2.0 --username admin --password admin --timeout "180s" check-vulnerabilities --project team1 --repository harbor --artifact build-PR-1 --force-scan --severity CRITICAL
```

You need to set following parameter:

- **--project**: The project name
- **--registry**: The registry name
- **--artifact**: The artifact name.
- **--severity**: The maximum severity before to exit with error (LOW, MEDIUM, HIGH or CRITICAL). Default it not check severity.
- **--force-scan**: To launch new scan and wait before to access on vulnerabilities reports. Default to `false`.

It return the following code:
- `0`: all work fine and severity is not to bad
- `1`: some internale errors or severity is to bad.

### Delete artifact (docker image)

It permit to delete a specific artifact (docker image).

We use it on Jenkins CI step just after destroy volatile environment for user tests on PR pipeline.

Sample of command:

```bash
harbor-cli --url https://harbor.company.com/api/v2.0 --username admin --password admin delete-artifact --project team1 --repository harbor --artifact build-PR-1
```

You need to set following parameter:

- **--project**: The project name
- **--registry**: The registry name
- **--artifact**: The artifact name.

It return the following code:
- `0`: all work fine
- `1`: some internale errors

## Contributing

You PR are always welcome. Please use the `main` branch as PR target.
Don't forget to add test if you add some functionalities.

To build, you can use the following command line:

```sh
make build
```

To lauch golang test, you can use the folowing command line:

```sh
make test
```

## Dev notes

### Generate mocks
```
mockgen -destination=harbor/mocks/mock_api.go -package=mocks github.com/disaster37/harbor-cli/harbor/api API,ArtifactAPI
```