# templ-iconify


[![release](https://img.shields.io/github/release/ksckaan1/templ-iconify.svg)](https://github.com/ksckaan1/templ-iconify/releases)
![Go Version](https://img.shields.io/badge/Go-1.24.2-%23007d9c)
[![Go report](https://goreportcard.com/badge/github.com/ksckaan1/templ-iconify)](https://goreportcard.com/report/github.com/ksckaan1/templ-iconify)
![coverage](https://img.shields.io/badge/coverage-none-green?style=flat)
[![Contributors](https://img.shields.io/github/contributors/ksckaan1/templ-iconify)](https://github.com/ksckaan1/templ-iconify/graphs/contributors)
[![LICENSE](https://img.shields.io/badge/LICENCE-MIT-orange?style=flat)](./LICENSE)

This CLI tool download icons over `Iconify` and converts them to `templ` files.

## Features

- Download icons over `Iconify`
- Convert icons to `templ` files
- Download multiple icons
- Download icons using wildcard
- Offline Usage Support as Default
- Parallel Download Support with workers

## Installation

```sh
go install github.com/ksckaan1/templ-iconify@latest
```

## Usage

```sh
templ-iconify <icon-name> -o <output-dir>
```

### CLI Examples:
- Download `mdi:home` icon to `./icons` directory (default: `./icons`)
  ```sh
  templ-iconify "mdi:home"
  ```

- Download icons starts with `mdi:home-` to `./icons` directory
  ```sh
  templ-iconify "mdi:home-*" -o ./icons
  ```

- Download icons satisfies `mdi:*` to `./icons` directory
  ```sh
  templ-iconify "mdi:*" -o ./icons
  ```

- Download with multiple expressions
  ```sh
  templ-iconify "mdi:home*" "solar:home*" -o ./icons
  ```

- Download all icons
  ```sh
  templ-iconify "*:*" -o ./icons
  ```

- Download with custom worker count (default: 10)
  ```sh
  templ-iconify "mdi:home" -o ./icons -w 20
  ```
  
### templ Example:

```templ
package templates

import "<module-name>/icons/mdi"

templ Page() {
  @mdi.Home(mdi.HomeProps{
    Height: "100px",
    Width: "100px",
    Color: "red", // default: "currentColor"
  })
}

```

## License

MIT License