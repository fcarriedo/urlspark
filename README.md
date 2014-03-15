# urlspark

[![Build Status](https://travis-ci.org/fcarriedo/urlspark.png?branch=master)](https://travis-ci.org/fcarriedo/urlspark)

Quick and dirty URL shortener that creates volatile (read *short lived*)
shortened URLs.

### Installing it

You need [golang](http://golang.org/) installed with the GOPATH set.

Just get it!

```
  $ go get github.com/fcarriedo/urlspark
```

### Running it

```
Usage:  urlspark [options]

Starts a URL shortener server

  -p=80: http port to run
  -redis="": redis address 'host:port'
  -ttl=60: The expiration time [time to live seconds]
```

### Using it

To generate a shortened short lived URL just POST a long URL to the base.
You should get the shortened URL as a response.

```
  $ curl --data-urlencode "url=http://ahost.com/and/very/long/url?param=some&key=val" http://urlsparkhost
  ...
  http://urlsparkhost/pK7z
```

Now point your browser to the generated URL: `http://urlsparkhost/pK7z`. It
should redirect you to `http://ahost.com/and/very/long/url?param=some&key=val`

The shortened URL should exist for as long as `exp` seconds (*60 sec by
default*)

If by any means you need to delete the generated shortened URL before it
expires:

```
  $ curl -X DELETE http://urlsparkhost/pK7z
```

...and the resource should cease to exist.

### Dependencies

The only dependency is on [redis](http://redis.io/) which works as its
datastore. Redis seems to make the most sense due to the nature of the app
(key, value, expiration) in addition to its blazing speed. Other datastores can
be implemented but expiration might need to be added if not an inherent
property of the backing datastore.

The service expects the redis server to be running and reachable from the
server to be able to perform its main operations.

### TODO

  * Think about implementing **single use URLs**

### License

urlspark is available under the [Apache License, Version
2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
