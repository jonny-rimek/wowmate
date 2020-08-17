#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { Wowmate } from '../lib/wowmate-stack';

/*
TODO: add 
interface EnvProps {
  prod: boolean;
}
update constructors 
constructor(scope: Construct, id: string, props?: EnvProps) {
	*/

const app = new cdk.App();
new Wowmate(app, 'wm', {
	env: {region: "us-east-1"}
});