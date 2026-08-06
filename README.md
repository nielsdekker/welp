# Welp

`Welp` is a crawler that enumerates string-like values on webpages. With
frontend frameworks a lot of API data or tokens could be hidden somewhere in the
JavaScript bundle and `welp` will attempt to find these.

> NOTE: Welp results in a lot of requests to the given target. Use only when you
> have permission

# How to use

The simplest call is as follows:

```bash
welp -u http://target.test
```

A JSON representation of the output can also be written to a file:

```bash
welp -u http://target.test --output out.json
```

## Adding prefixes

It's possible the JavaScript code contains a snippet as follows:

```js
const apiVersion = "/rest/v2/"
fetch(apiVersion + "users")
```

`Welp` will see two strings values and checks the `"/rest/v2/"` and `"/users"`
endpoints instead of `"/rest/v2/users"`. For these situations a `--prefix` flag
can be passed 

```bash
welp -u http://target.test --prefix /rest/v2/ --prefix /rest/v3/
```

This will result in each found value to also be tested with the given prefixes.
In short the above example will make calls to:

- `/rest/v2/`
- `/rest/v2/rest/v2/`
- `/rest/v3/rest/v2/`
- `/users`
- `/rest/v2/users`
- `/rest/v3/users`

## Using modules

Some modules are included to parse the output. These are:

- `text`, Will also print all the found text values
- `token`, Matches certain token types
- `entropy`, Checks the text has a certain entropy level. In case `tokens`
  doesn't find it

Modules can be used as follows:

```bash
welp -u http://target.test --module text -m token
```

## Filtering the output

By default 404 pages are filtered in the output but additional filter options
can be given.

```bash
# Removes all responses with a 404 or 500 status code from the output
welp -u http://target.test -fc 404 --filter-code 500

# Removes all responses with a content type that matches # `-ft text` matches
# `text/html`, `test/css`, etc.
welp -u http://target.test -ft text -ft image
```
