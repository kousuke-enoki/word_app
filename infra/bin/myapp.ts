#!/usr/bin/env node
import 'source-map-support/register';
import { App } from 'aws-cdk-lib';
import { NetworkStack } from '../lib/network-stack';
import { DbStack } from '../lib/db-stack';
import { AppStack } from '../lib/app-stack';
import { OpsStack } from '../lib/ops-stack';

const app  = new App();
const env  = { account: process.env.CDK_DEFAULT_ACCOUNT,
               region : process.env.CDK_DEFAULT_REGION };

const net = new NetworkStack(app,'NetStack',{ env });
const db  = new DbStack(app,'DbStack',{ env, vpc: net.vpc });
// AppStack は VPC と DB を参照する
// そのため、net と db の後に定義する必要がある
new AppStack(app,'AppStack',{
  env,
  vpc: net.vpc,
  secret: db.secret,
  db: db.db,
  lambdaSg: db.lambdaToDbSecurityGroup,
});

new OpsStack(app,'OpsStack',{
  env,
  db: db.db,               // ← RDS だけ渡す
  stopHourJst: 0,          // 好みで変更可
  startHourJst: 9,
});
