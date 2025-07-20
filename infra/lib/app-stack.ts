import { Stack, StackProps, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda  from 'aws-cdk-lib/aws-lambda';
import * as logs    from 'aws-cdk-lib/aws-logs';
import * as apigw   from 'aws-cdk-lib/aws-apigateway';
import * as ec2     from 'aws-cdk-lib/aws-ec2';
import * as ecr     from 'aws-cdk-lib/aws-ecr';
import * as secrets from 'aws-cdk-lib/aws-secretsmanager';
import * as rds     from 'aws-cdk-lib/aws-rds';

export interface AppStackProps extends StackProps {
  vpc: ec2.IVpc;
  secret: secrets.ISecret;
  db: rds.IDatabaseInstance;
  lambdaSg: ec2.ISecurityGroup;
}
/**
 * アプリケーションスタック
 * ECR イメージを Lambda で実行し、API Gateway で公開する
 */
export class AppStack extends Stack {
  constructor(scope: Construct, id: string, { vpc, secret, db, lambdaSg, ...rest }: AppStackProps) {
    super(scope, id, rest);

    // ECR イメージ → Lambda
    const repo = ecr.Repository.fromRepositoryName(this,'Repo','wordapp-backend');
    const fn = new lambda.DockerImageFunction(this,'ApiFn',{
      code: lambda.DockerImageCode.fromEcr(repo,{ tagOrDigest: 'latest' }),
      vpc,
      vpcSubnets: { subnetType: ec2.SubnetType.PRIVATE_ISOLATED },
      securityGroups: [lambdaSg],
      memorySize: 256,
      timeout: Duration.seconds(10),
      logRetention: logs.RetentionDays.THREE_DAYS,
      environment:{
        DB_SECRET_ARN: secret.secretArn,
        DB_NAME: 'postgres',
      },
    });

    /* 接続許可 */
    secret.grantRead(fn);          // 認証情報を読む権限
    // db.connections.allowDefaultPortFrom(fn); // SG で 5432 許可
    // // もしくは db.grantConnect(fn); でも可

    new apigw.LambdaRestApi(this,'Api',{ handler: fn });
  }
}
