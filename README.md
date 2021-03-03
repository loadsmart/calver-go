# calver-go - The Calendar Versioner for go

Parse and increment version based on [calver.org](https://calver.org/) strategy.

[![CircleCI](https://circleci.com/gh/loadsmart/calver-go.svg?style=svg&circle-token=66bca7a1bfc187f8e93d8876f70596be9eeff346)](https://circleci.com/gh/loadsmart/calvergo)

Install
```shell
go get github.com/loadsmart/calver-go
```

Usage
```go
import "github.com/loadsmart/calver-go/calver"

latest, err := calver.Parse("YYYY.MM.DD.MICRO", "2019.11.27.1")
micro := latest.Next() // if today is 11/27/2019 then version == 2019.11.27.2
major := latest.Next() // if today is 11/28/2019 then version == 2019.11.28.1

brandNewVersion := calver.NewVersion('YYYY.MM.DD.MICRO', 0) // if today is 11/28/2019 then version == 2019.11.28.1
brandNewVersion.String() // "2019.11.28.1"
```

### Caveats
* support only the conventions below:
  * YYYY
  * YY
  * 0Y
  * MM
  * M0
  * 0M
  * DD
  * D0
  * 0D
  * MICRO (a.k.a build)

TODO
* export as command-line tool
* add support to other conventions


## License

[MIT](./LICENSE)
