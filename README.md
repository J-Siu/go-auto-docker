# go-auto-docker

Automate update single Alpine package docker container. Update dockerfile, change log, build test, commit, git tag according to package version.

- [Install](#install)
- [Usage](#usage)
- [Limitation](#limitation)
- [License](#license)

<!--more-->
### Install

Go install

```sh
go install github.com/J-Siu/go-auto-docker@latest
```

Download

- https://github.com/J-Siu/go-auto-docker/releases

### Usage

```sh
Automate update for README.md change log, apply tag according to package version. Also handle test build, git commit.

Usage:
  go-auto-docker [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Print configurations
  db          DB commands
  help        Help about any command
  update      Update Alpine package version

Flags:
      --config string   config file (default "~/.config/go-auto-docker.json")
  -d, --debug           enable debug
  -h, --help            help for go-auto-docker
  -v, --verbose         enable debug
      --version         version for go-auto-docker

Use "go-auto-docker [command] --help" for more information about a command.
```

Update:

```sh
Update Alpine package version

Usage:
  go-auto-docker update [flags]

Aliases:
  update, u

Flags:
  -b, --buildTest   so not perform docker build
  -c, --commit      apply git commit. Only work with -save
  -h, --help        help for update
  -s, --save        write back to project folder (cancel on error)
  -t, --tag         apply git tag. (only work with --commit)
  -u, --updateDb    update Alpine package database

Global Flags:
      --config string   config file (default "~/.config/go-auto-docker.json")
  -d, --debug           enable debug
  -v, --verbose         enable debug
```

```sh
go-auto-docker update \
--buildTest \ # Test docker build
--commit \    # Git commit
--save \      # Save back to original folder
--tag \       # Git tag with new version
docker_*      # Handle multiple repository directories
```

### Limitation

- Assume single Alpine package docker container
- Does not work in MacOS
- Dockerfile
  - "LABEL version:" equal to package version
  - `RUN` install line should specify version
  - Assume `main` and `community` repository
  - Detect `testing` branch via `edge/testing`

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
