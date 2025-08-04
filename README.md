# tp
**🗃️ `tp` is a CLI tool for conveniently generating software license and gitignore templates for your projects ⚒️.** No longer do you need to manually scour the internet in order to find the correct `LICENSE` and `.gitignore` files, `tp` will automatically go out and fetch your desired `LICENSE` templates from <https://github.com/spdx/license-list-data> and `.gitignore` templates <https://github.com/github/gitignore>.

# Install
- Dependencies
	- [go](https://go.dev)

```shell
go install github.com/JessebotX/tp@latest
```

# Basic Usage
To see general help/usage information of `tp`, run:

```shell
tp --help
```

You can see help/usage info for specific commands as well.

```shell
tp --help license
tp --help gitignore
```

## LICENSE
To see all available `LICENSE` templates (SPDX license identifiers), run:

```shell
tp license -l
```

To fetch a `LICENSE` template, run (NOTE: case matters i.e. `MIT` is not the same as `mit`):

```shell
# format
tp license <SPDX-License-Identifiers...>

# e.g. download the "MIT" license data as a file named LICENSE in your current directory
tp license MIT

# e.g. download the "Apache 2" license data into "/path/to/project/LICENSE.txt"
tp license Apache-2.0 -o /path/to/project/LICENSE.txt

# e.g. print the concatenation of the "MIT" and the "GNU GPL 3 only" to your terminal (i.e. to stdout)
tp license MIT GPL-3.0-only --stdout
```

## .gitignore
To see all available `.gitignore` templates, run:

```shell
tp gitignore -l
```

Generating a `.gitignore` template is similar to generating a `LICENSE` template
```shell
# format
tp gitignore <SPDX-License-Identifiers...>

# e.g. download the Go .gitignore template as a file named ".gitignore" in your current directory
tp gitignore Go

# e.g. download the Node.js gitignore template as "/path/to/project/.gitignore"
tp gitignore Node -o /path/to/project/.gitignore

# e.g. print the concatenation of Go programming language's gitignore template and the Hugo
#      static site generator's gitignore template to your terminal (i.e. to stdout)
tp gitignore Go community/Golang/Hugo --stdout
```

# License / Permissions
See [LICENSE](./LICENSE) for more information.
