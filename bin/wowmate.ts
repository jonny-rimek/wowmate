#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { Wowmate } from '../lib/wowmate-stack';

const app = new cdk.App();
new Wowmate(app, 'wm', {
	env: {region: "us-east-1", account: "302123354508"},
	hostedZoneId: "Z08580822XS57UHUUVCD4",
	hostedZoneName: "wowmate.io",
	domainName: "wowmate.io",
	apiDomainName: "api.wowmate.io",
});

new Wowmate(app, 'wm-dev', {
	env: {region: "us-east-1", account: "461497339039"},
	hostedZoneId: "Z09026202SZR8MRVSF1BQ",
	hostedZoneName: "dev.wowmate.io",
	domainName: "dev.wowmate.io",
	apiDomainName: "api.dev.wowmate.io",
});

new Wowmate(app, 'wm-preprod', {
	env: {region: "us-east-1", account: "500489575211"},
	hostedZoneId: "Z09026202SZR8MRVSF1BQ",
	hostedZoneName: "preprod.wowmate.io",
	domainName: "preprod.wowmate.io",
	apiDomainName: "api.preprod.wowmate.io",
});
