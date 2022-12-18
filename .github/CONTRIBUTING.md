# CONTRIBUTING

Any PullRequest for the better is welcome!

- Branch to PR: `main`
- Issues:
  - Please include a simple code snippet to reproduce the issue. It will help us a lot to fix the issue.
  - Issues can be in Japanese and Spanish rather than English if you prefer.

## Tests and CIs

You need to pass the below before review.

- `go test -race ./...`
  - On Go 1.18 to the latest.
- `golangci-lint run`
  - See the `.golangci.yml` for the configuration. Requires `golangci-lint` v1.50.1 or later.
- `golint ./...`
  - `golint` is deprecated though we still use it to find missing comments.
- `go test -cover ./...`
  - Please keep the code coverage up to 100%.
  - We recommend to use [go-carpet](https://github.com/msoap/go-carpet) to find which lines are not covered.
    - `go-carpet -mincov 99.9 ./...`

We have CIs to check these. So we recommend to [draft PR](https://github.blog/2019-02-14-introducing-draft-pull-requests/) before you implement something.

For convenience, there is a `docker-compose.yml` for the above tests.

```bash
# Run test on Go 1.18
docker compose run v1_18

# Run test on Go 1.19
docker compose run v1_19

# Run test on latest Go
docker compose run latest

# Run linters (golangci-lint)
docker compose run lint
```
