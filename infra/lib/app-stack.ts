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
  natEnabled: boolean; // ネットワーク設定を渡す
}
/**
 * アプリケーションスタック
 * ECR イメージを Lambda で実行し、API Gateway で公開する
 */
export class AppStack extends Stack {
  constructor(scope: Construct, id: string, { vpc, secret, db, lambdaSg, natEnabled,...rest }: AppStackProps) {
    super(scope, id, rest);

    const appSecret = new secrets.Secret(this, 'AppSecret', {
      secretName: 'wordapp/app-secrets',
      generateSecretString: {
        secretStringTemplate: JSON.stringify({
          LINE_CLIENT_ID    : 'change-me',
          LINE_CLIENT_SECRET: 'change-me',
          LINE_REDIRECT_URI : 'https://example.com/callback',
          JWT_SECRET        : 'change-me',
        }),
        generateStringKey: 'placeholder', // 何か1つは自動生成、未使用でもOK
      },
    });

    const subnets = natEnabled
      ? { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS }
      : { subnetType: ec2.SubnetType.PRIVATE_ISOLATED };

    const imageTagOrDigest =
      this.node.tryGetContext('imageDigest') ??
      this.node.tryGetContext('imageTag') ?? 'lambda';


    // ECR イメージ → Lambda
    const repo = ecr.Repository.fromRepositoryName(this,'Repo','wordapp-backend');
    const fn = new lambda.DockerImageFunction(this,'ApiFn',{
      code: lambda.DockerImageCode.fromEcr(repo,{ tagOrDigest: imageTagOrDigest }),
      vpc,
      vpcSubnets: subnets,
      securityGroups: [lambdaSg],
      memorySize: 256,
      timeout: Duration.seconds(30), // ← 少し余裕を持たせる
      logRetention: logs.RetentionDays.THREE_DAYS,
      environment:{
        APP_ENV: 'production',
        DB_HOST: db.instanceEndpoint.hostname,
        DB_PORT: db.instanceEndpoint.port.toString(),
        DB_NAME: 'postgres',
        DB_SECRET_ARN: secret.secretArn,
        APP_SECRET_ARN: appSecret.secretArn, // ←ここは将来 既存SecretのARNに差し替える
      },
    });

    /* 接続許可 */
    secret.grantRead(fn);          // 認証情報を読む権限
    appSecret.grantRead(fn);
    // db.connections.allowDefaultPortFrom(fn); // SG で 5432 許可
    // // もしくは db.grantConnect(fn); でも可

    new apigw.LambdaRestApi(this,'Api',{ handler: fn });
  }
}
