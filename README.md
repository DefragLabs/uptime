# uptime

[![CircleCI](https://circleci.com/gh/DefragLabs/uptime/tree/master.svg?style=svg)](https://circleci.com/gh/DefragLabs/uptime/tree/master)

# Goals

- Server uptime monitoring
  - Frequency (seconds/minutes/hours/days)
  - url
  - http/https
  - headers

Api's internal, but if required should be done easily.

## Tests
Run tests using this command `go test ./...`. 

Currently run them from the api container.

## Development

We add `structs` (From structs external package) tags in go struct so the output keys are properly formatted.
This package doesn't honour `json` tags.

## Monitor URL's

Monitor url's configuration matrix

| Frequency | Unit   |
|-----------|--------|
| 30        | second |
| 1         | minute |
| 5         | minute |
| 15        | minute |
| 30        | minute |
