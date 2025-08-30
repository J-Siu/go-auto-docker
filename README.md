### Limitation

- Does not work in MacOS
- Dockerfile
  - "LABEL version:" equal to package version
  - `RUN` install line should specify version
  - Assume `main` and `community` branch
  - Detect `testing` branch via `edge/testing`
