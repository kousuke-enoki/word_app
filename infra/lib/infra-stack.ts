import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as lambda from 'aws-cdk-lib/aws-lambda';

export class InfraStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const vpc = new ec2.Vpc(this, 'Vpc', { maxAzs: 2 });

    new lambda.Function(this, 'Hello', {
      runtime: lambda.Runtime.NODEJS_18_X,
      code: lambda.Code.fromInline('exports.handler=()=>({statusCode:200,body:"hi"})'),
      handler: 'index.handler',
      vpc,
    });
  }
}
