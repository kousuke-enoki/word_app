import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as ecr from 'aws-cdk-lib/aws-ecr';
import * as secrets from 'aws-cdk-lib/aws-secretsmanager';

export interface AppStackProps extends StackProps {
  vpc: ec2.IVpc;
  secret: secrets.ISecret;
}

export class AppStack extends Stack {
  constructor(scope: Construct, id: string, { vpc, secret, ...rest }: AppStackProps) {
    super(scope, id, rest);

    const repo = ecr.Repository.fromRepositoryName(
      this,
      'BackendRepo',
      'wordapp-backend',
    );

    const fn = new lambda.DockerImageFunction(this, 'ApiFn', {
      code: lambda.DockerImageCode.fromEcr(repo, { tag: 'latest' }),
      vpc,
      memorySize: 256,
      environment: {
        DB_SECRET_ARN: secret.secretArn,
        DB_NAME: 'postgres',
      },
    });

    secret.grantRead(fn);

    new apigw.LambdaRestApi(this, 'Api', { handler: fn });
  }
}
