# uptime

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

## Monitor URL's

Monitor url's configuration matrix

| Frequency | Unit   |
|-----------|--------|
| 30        | second |
| 1         | minute |
| 5         | minute |
| 15        | minute |
| 30        | minute |
