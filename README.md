# wowmate mono repo

## manual steps:

- creating the public hosted zone in route53
- add s3 import role
- enable http api logs
- enable cloudfront advanced metrics

### update cdk cli

`sudo npm i -g aws-cdk`

#### install psql on the bastion station 

`sudo amazon-linux-extras install postgresql11`

the bastion will be replaced when a new AMI version is released