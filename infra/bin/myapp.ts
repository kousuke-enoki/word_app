#!/usr/bin/env node
import 'source-map-support/register';
import { App } from 'aws-cdk-lib';
import { NetworkStack } from '../lib/network-stack';
import { DbStack } from '../lib/db-stack';
import { AppStack } from '../lib/app-stack';

const app  = new App();
const env  = { account: process.env.CDK_DEFAULT_ACCOUNT,
               region : process.env.CDK_DEFAULT_REGION };

const net  = new NetworkStack(app, 'NetStack', { env });
const db   = new DbStack(app,  'DbStack',  { env, vpc: net.vpc });
new AppStack(app,'AppStack',  { env, vpc: net.vpc, secret: db.secret });
