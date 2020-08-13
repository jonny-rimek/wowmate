#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { WowmatePipelineStack } from '../lib/wowmate-stack';

/*
TODO: add 
interface EnvProps {
  prod: boolean;
}
update constructors 
constructor(scope: Construct, id: string, props?: EnvProps) {
	*/

const app = new cdk.App();
new WowmatePipelineStack(app, 'WoWM', {
	env: {region: "us-east-1"}
});