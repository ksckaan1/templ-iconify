# templ-iconify

This CLI tool download icons over `Iconify` and converts them to `templ` files.

## Installation

```sh
go install github.com/ksckaan1/templ-iconify@latest
```

## Usage

```sh
templ-iconify <icon-name> -o <output-dir>
```

### Example:
- Download `mdi:home` icon to `./icons` directory
  ```sh
  templ-iconify "mdi:home" -o ./icons
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
  



