import { Stack, StackProps, RemovalPolicy } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as rds from 'aws-cdk-lib/aws-rds';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as secrets from 'aws-cdk-lib/aws-secretsmanager';

export interface DbStackProps extends StackProps {
  vpc: ec2.IVpc;
}

export class DbStack extends Stack {
  public readonly secret: secrets.ISecret;
  public readonly db: rds.DatabaseInstance;

  constructor(scope: Construct, id: string, { vpc, ...rest }: DbStackProps) {
    super(scope, id, rest);

    this.secret = new rds.DatabaseSecret(this, 'DbSecret', {
      username: 'postgres',
    });

    this.db = new rds.DatabaseInstance(this, 'Postgres', {
      engine: rds.DatabaseInstanceEngine.postgres({
        version: rds.PostgresEngineVersion.VER_15,
      }),
      vpc,
      credentials: rds.Credentials.fromSecret(this.secret),
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T4G,
        ec2.InstanceSize.MICRO,
      ),
      allocatedStorage: 20,
      removalPolicy: RemovalPolicy.DESTROY,   // 検証用
    });
  }
}
