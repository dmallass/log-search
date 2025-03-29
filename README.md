### LOG COLLECTION

## INSTALLATION AND SETUP

```
git clone git@github.com:dmallass/log-search.git
cd log-search
```

### Install Dependencies

Make sure you have Go (>=1.21) and ripgrep installed:
```
[Install Go>=1.21](https://go.dev/doc/install)
brew install ripgrep
go get -u github.com/gin-gonic/gin
go mod tidy
```

### Run the API

```
go build -o .
./cribl_take_home
```

## DOCKER SETUP

```
docker build -t search-api .

docker run -p 8080:8080 search-api:latest
```

### Test the API

#### Required Parameters

| Name | Value |
| --- | --- |
| logfile    |  /var/log/app.log |
|  q         |  ERROR       |


#### Optional Parameters 

| Name | Value  |
| --- | --- |
| page         |  1         |
| limit        |  5         |
| searchMode   | fulltext or regex |

#### full text search
```
curl "http://localhost:8080/search?page=1&limit=5&logfile=/var/log/app.log&q=v1/shoppers"
```
#### regex search
```
curl "http://localhost:8080/search?page=1&limit=5&logfile=/var/log/app.log&q=.*/v1/candidates&searchMode=regex"

```

## Performance Optimization

Adding additional ripgrep (rg) flags could make the performance even better. With increseing log file size or complexity, considering the current performance of ripgrep for 1GB file, we could scale up the server resources cpu and memory.

Another optimization is to add a cache. For an optimized performane for repeated searches, The results could be cached at the broker application or at the CDN level. The response headers from the search application could include the logfile last-updated timestamp . The cache will return results matched with the lastupdated timestamp key. 

## Known Issues

#### Performance BenchMark

| Method of packaging | Search time in seconds | 
| ---- | -----| 
| binary executable | 0.8ms - 2ms |
| Containerized application (Docker) | upto 4s |

When run in a container, virtualization is affecting the search performance. When run in apple silicon, the search time is under 2 seconds.

#### Authentication

REST api endpoint doesn't check if the user requesting the search has the correct privilages. 

1. JWT Tokens for authentication
```Authentication``` header will include ```<Bearer <JWTToken>> ``` . Application can have a authtication middleware that will validate the token before allowing to access the endpoint. if token is missing or invalid, a 401 response will be returned. 

2. Check user's privilages for directory access. make sure the user requesting for logs has the correct privilages to access the directory. 
