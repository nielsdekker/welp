# Welp

The goal of `welp` is to find URL's and other string-like content on webpages.
For example a _Single Page Application_ could have a goldmine of API URL's
somewhere in the javascript, `welp` will try to find these. Another possibility
is tokens or other secrets that were accidentally exposed in the frontend code.

# How to use

```bash
welp -u http://target.test
```

# Technical

What `welp` does is find content within `"`, `'`, or `\``. Then makes a request
for all the values that match and repeats the same process. For example:

```html
<html>
    <head>
        <script type="javascript" src="/script.js"></script>
    <head>
    <body>
        <span>My awesome webpage</span>
    </body>
</html>
```

> Note, the text "My awesome webpage" is not in quotes and therefore this is
> skipped.

`welp` Will find the following:

- `javascript`, which results in the call `GET /javascript`
- `/script.js`, which triggers a call to `GET /script.js`

On those pages the same logic is repeated.

## Safeguards

There are some safeguards in place to prevent endless loops.

- `40x` and `50x` Responses will be reported as a result but content on these
  do not trigger another crawl
- A md5 sum is calculated, when a page matches another page the result is
  reported but no additional requests will be made
    - Reasoning is _single page applications_ could return the exact same
      `index.html` for each non-api endpoint. The routing is determined in
      the JavaScript itself.
- If the `path` gets too long, this is mostly to filter out large `eval("...")`
  blocks that could be present in the JavaScript
- When the new URL targets another domain it will not be crawled
