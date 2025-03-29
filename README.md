### LOG COLLECTION

## INSTALLATION AND SETUP

```
git clone 
cd cribl_take_home
```

## Install dependencies 

Make sure you have Go (>=1.21) and ripgrep installed:

```
go mod tidy

```

## Run the API

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

```

| Required Parameters       |
| _________________________ |
| Param name | Param value  |
| ________   | ___________  |
| logfile    |  /var/log/app.log     |
|  q         |  ERROR       |


| Optional Parameters       |
| _________________________ |
| Param name | Param value  |
| ________   | __________   |
| page         |  1         |
| limit        |  5         |
| searchMode   | fulltext or regex |

```
#### full text search
```
curl "http://localhost:8080/search?page=1&limit=5&logfile=/var/log/app.log&q=v1/shoppers"
```
#### regex search
```
curl "http://localhost:8080/search?page=1&limit=5&logfile=/var/log/app.log&q=.*/v1/candidates&searchMode=regex"

```

## Performance Optimization

The search performance is determined by the ripgrep performance. 
Adding additional ripgrep (rg) flags could make the performance even better. With an incresed log file size or complexity, considering the current performance of ripgrep for 1GB file, we could scale up the server resources cpu and memory.

Another optimization is to add a cache. For an optimized performane for repeated searches, The results could be cached at the broker application or at the CDN level. The response headers from the search application could include the logfile created timestamp and last-updated timestamp . The cache will return results matched with the lastupdated timestamp key. 


