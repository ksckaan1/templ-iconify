# templ-iconify

This CLI tool download icons over `Iconify` and converts them to `templ` files.

## Features

- Download icons over `Iconify`
- Convert icons to `templ` files
- Download multiple icons
- Download icons using wildcard
- Offline Usage Support as Default

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