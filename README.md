# wowmate mono repo

## manual steps:

- creating the public hosted zone in route53
- enable http api logs
- enable cloudfront advanced metrics

#### install psql on the bastion station 

`sudo amazon-linux-extras install postgresql11`

the bastion will be replaced when a new AMI version is released

#### go mod

e.g. `go mod init github.com/jonny-rimek/wowmate/services/api/combatlogs-combatlog-uuid-damage`

### go mod replace

`go mod edit -replace y=github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize./normalize`

add code that uses the package e.g. `normalize.Normalize`

manually import the package in the import section by adding `"github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize"`