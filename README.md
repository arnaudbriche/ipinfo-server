# ipinfo-server

## Build Docker image

```sh
docker \
    buildx \
    build \
    --platform linux/amd64 \
    -t arnaudbriche/ipinfo-server \
    --push \
    .
```

## Run

```
docker run -p 8081:8080 -it arnaudbriche/ipinfo-server
```

## Usage

### Getting egress IP info

```sh
curl -v 'http://localhost:8081'
```

### Resolving hostname

```sh
curl -v 'http://localhost:8080/resolve?hostname=google.com'
```