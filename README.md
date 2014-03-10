# urlspark

Quick and dirty URL shortener that creates volatile (read *short lived*)
shortened URLs.

### Dependencies

The only dependency is on [redis](http://redis.io/) which works as its
datastore.

The service expects the redis server to be running and reachable to be able to
perform its main operations.

### Installing it

You need [golang](http://golang.org/) installed with the GOPATH set.

Just get it!

```
  $ go get github.com/fcarriedo/urlspark
```

### Running it

Running:

```
  $ urlspark
```

The previous runs the service with the following defaults:

  * `p := 80` - Runs the service on port **80**
  * `exp := 60` - Every generated short URL expires in **60 seconds**
  * `redis := localhost:6379` - Redis server expected on **localhost:6379**

You can customize any of the previous parameters. A fully customized startup
would look like the following:

```
  $ urlspark -p 8080 -exp 30 -redis "myredishost:6399"
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

### License

urlspark is available under the [Apache License, Version
2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
