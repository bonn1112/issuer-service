# issuing-service

## Testing

Build the docker image for running testing in container

```bash
docker build -f ./Dockerfile.test -t issuing-service-test .
docker run --rm issuing-service-test
```