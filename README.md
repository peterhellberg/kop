# Köp 📝

[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/peterhellberg/kop#license-mit)

A small example of using [oto](https://github.com/pacedotdev/oto) for RPC code generation.

## Installation

```
go install github.com/peterhellberg/kop/cmd/...@latest
```

## Usage

Run `kop-server` in one terminal (this will by default start a web server on `localhost:12432`)

> [!NOTE]
> If you change the `PORT` for `kop-server`, then you will have to set the `KOP_ENDPOINT` variable accordingly.
>
> (Default is `http://localhost:12432/rpc/`)

Then run `kop Eggs Milk Flour` to create an initial list using the command line interface.

```md
 - EGGS
 - FLOUR
 - MILK
```

> [!IMPORTANT]
> If you speak Swedish, then you will want to `alias köp='kop'` (and if not, then you might want to `alias buy='kop'`)

You can then remove an item from the list by calling `kop no Eggs`

```md
 - FLOUR
 - MILK
```

If you do not want to use the cli then you can use [cURL](https://curl.se) directly

```sh
curl -d '{}' http://localhost:12432/rpc/List.Items
```
```json
{
  "items": [
    "FLOUR",
    "MILK"
  ]
}
```

And even add something else to the list 🍺

```sh
curl -d '{"items": ["Beer"]}' http://localhost:12432/rpc/List.Add
```
```json
{
  "items": [
    "BEER",
    "FLOUR",
    "MILK"
  ]
}
```

## Definitions

The definitions for the RPC service are found in [definitions/definitions.go](definitions/definitions.go)

<img src="https://assets.c7.se/svg/viking-gopher.svg" align="right" width="30%" height="300">

## License (MIT)

Copyright (c) 2023 [Peter Hellberg](https://c7.se)

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
>
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
