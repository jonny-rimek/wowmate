module github.com/jonny-rimek/wowmate/services/api/combatlogs/advanced-damage/_combatlog_uuid/get

go 1.16

replace github.com/jonny-rimek/wowmate/services/common/golib => ./../../../../../../common/golib

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/aws/aws-sdk-go v1.38.36
	github.com/aws/aws-xray-sdk-go v1.4.0
	github.com/jonny-rimek/wowmate/services/common/golib v0.0.0-20210508172143-e8e5a2a594df
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
