# cachanais

Populate cache by crawling pages.

## Why

When purging cache (or restaring), you may want to repopulate cache.
_cachanais_ crawl pages by following `href`.

## Usage

* basic

``` shell
$ cachanais --url https://my-website.com
Visiting https://my-website.com/
Visiting https://my-website.com/page1
Visiting https://my-website.com/page2
```

* with an ssl offloader

``` shell
$ cachanais --url https://my-website.com --address http://localhost:8080
Visiting https://localhost:8080/
Visiting https://localhost:8080/page1
Visiting https://localhost:8080/page2
```

* with custom headers and cookies

``` shell
$ cachanais --url https://my-website.com --address http://localhost:8080 --cookies mycookie:value --headers X-My-Header:value
Visiting https://localhost:8080/
Visiting https://localhost:8080/page1
Visiting https://localhost:8080/page2
```

## Name

_cachanais_ is the gentile of the city of _Cachan_.
Because we populate cache.
