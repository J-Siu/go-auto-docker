# go-auto-docker

### Table Of Content
<!-- TOC -->
<!-- /TOC -->
### Limitation

- Assume single Alpine package docker container
- Does not work in MacOS
- Dockerfile
  - "LABEL version:" equal to package version
  - `RUN` install line should specify version
  - Assume `main` and `community` repository
  - Detect `testing` branch via `edge/testing`

### Change Log

- v0.5.0
  - Feature completed
- v0.5.1
  - update to go-helper/v2
- v0.5.2
  - `TypeReadme` use property
  - fix `TypeDocker` "package=version" line extraction
  - update go-helper/v2

### License

The MIT License (MIT)

Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
