# Log Search

## Installation and Setup

```
git clone git@github.com:dmallass/log-search.git
cd log-search
```

### Install Dependencies

Make sure you have Go (>=1.21) and ripgrep installed:

[Install Go>=1.21](https://go.dev/doc/install)
```
brew install ripgrep
go get -u github.com/gin-gonic/gin
go mod tidy
```

### Run the API

```
go build -o .
./cribl_take_home
```

## Docker Setup

```
docker build -t search-api .

docker run -p 8080:8080 search-api:latest
```

### Test the API

#### Required Parameters

| Name | Description | Value |
| --- | --- |--- |
| logfile    | file name to search for logs |  /var/log/app.log |
|  q         |  search query. You can provide basic keyword/text or regex. If regex is given, add searchMode=regex |   ERROR    or .*/v1/candidates   |
| searchMode   | fulltext or regex |

#### Optional Parameters 

| Name | Value  |
| --- | --- |
| page         |  1         |
| limit        |  5         |

#### full text search
```
curl "http://localhost:8080/search?page=1&limit=5&logfile=/var/log/app.log&q=v1/shoppers"
```
#### regex search
```
curl "http://localhost:8080/search?page=1&limit=5&logfile=/var/log/app.log&q=.*/v1/candidates&searchMode=regex"

```

## Performance Optimization

#### Performance BenchMark

| Method of packaging | File Size |Search time in seconds | Operating System | CPU | Memory
| ---- | -----| ---- | ----- | ----| ---- | 
| binary executable | 1GB |0.8s - 1.8s | MacOS M1 Pro (darwin arm64 arch) | 8 CPU cores | 16 GB
| Containerized application (Docker) | 1GB |2s - 3s | alpine linux | 8 CPU Cores | 14 GB
| binary executable | 10GB |20s - 38s | MacOS M1 Pro (darwin arm64 arch) | 8 CPU cores | 16 GB
| Containerized application (Docker) | 10GB | 2m25s  | alpine linux | 8 CPU Cores | 14 GB

Adding additional ripgrep (rg) flags could improve the performance. With increasing log file size or complexity, considering the current performance of ripgrep for 1GB file, we could scale up the server resources cpu and memory.

Another optimization is to add a cache. For an optimized performance for repeated searches, The results could be cached at the broker application or at the CDN level. The response headers from the search application could include the logfile last-updated timestamp . The cache will return results matched with the lastupdated timestamp key. 

## Known Issues

#### Authentication

REST api endpoint doesn't check if the user requesting the search has the correct privilages. 

1. JWT Token can be used for authentication.
```Authentication``` header will include ```<Bearer JWTToken> ``` . Application can have a authentication middleware that will validate the token before allowing to access the endpoint. If token is missing or invalid, a 401 response will be returned. 

2. Check user's privilages for directory access. Make sure the user requesting for logs has the correct privilages to access the directory. 
