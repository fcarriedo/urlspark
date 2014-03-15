# urlspark

[![Build Status](https://travis-ci.org/fcarriedo/urlspark.png?branch=master)](https://travis-ci.org/fcarriedo/urlspark)

Quick and dirty URL shortener that creates volatile (read *short lived*)
shortened URLs.

### Installing it

You need [golang](http://golang.org/) installed with the GOPATH set.

Just get it!

    $ go get github.com/fcarriedo/urlspark

### Running it

```
Usage:  urlspark [options]

Starts a URL shortener server

  -p=80: http port to run
  -redis="": redis address 'host:port'
  -ttl=60: The expiration time [time to live seconds]
```

By default, it doesn't have any infrastructure dependencies and uses an *in
memory* datastore. If `-redis=host:port` is specified then it will try to
connect to [redis](http://redis.io/) and use it its datastore. **A redis
datastore might be the best suited for production deployments.**

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

The shortened URL should exist for as long as `ttl` seconds (*60 sec by
default*)

If by any means you need to delete the generated shortened URL before it
expires:

    $ curl -X DELETE http://urlsparkhost/pK7z

...and the resource should cease to exist.

### License

urlspark is available under the [Apache License, Version
2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
