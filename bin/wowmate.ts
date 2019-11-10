#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { WowmateStack } from '../lib/wowmate-stack';

const app = new cdk.App();
new WowmateStack(app, 'WowmateStack');
