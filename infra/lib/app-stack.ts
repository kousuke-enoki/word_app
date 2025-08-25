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
  secret: secrets.ISecret;           // DB 用 Secret
  db: rds.IDatabaseInstance;
  lambdaSg: ec2.ISecurityGroup;
  natEnabled: boolean;
  appSecretArn?: string;             // APP 用 Secret の ARN を受け取れるようにする
  createSmVpce?: boolean;
}
/**
 * アプリケーションスタック
 * ECR イメージを Lambda で実行し、API Gateway で公開する
 */
export class AppStack extends Stack {
  constructor(scope: Construct, id: string, { vpc, secret: dbSecret, db, lambdaSg, natEnabled, appSecretArn, createSmVpce,...rest }: AppStackProps) {
    super(scope, id, rest);

    const appSecret = secrets.Secret.fromSecretCompleteArn(
      this,
      'AppSecretImported',
      appSecretArn ?? 'arn:aws:secretsmanager:ap-northeast-1:381492105871:secret:wordapp/app-k4I6ng',
    );

    const subnets = natEnabled
      ? { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS }
      : { subnetType: ec2.SubnetType.PRIVATE_ISOLATED };

    // ECR イメージ → Lambda
    const repo = ecr.Repository.fromRepositoryName(this,'Repo','wordapp-backend');
    const fn = new lambda.DockerImageFunction(this,'ApiFn',{
      code: lambda.DockerImageCode.fromEcr(repo,{ tagOrDigest: 'latest' }),
      vpc,
      vpcSubnets: subnets,
      securityGroups: [lambdaSg],
      memorySize: 256,
      timeout: Duration.seconds(30), // ← 少し余裕を持たせる
      logRetention: logs.RetentionDays.THREE_DAYS,
      environment:{
        APP_ENV: 'production',
        GIN_MODE: 'release',
        APP_PORT: '8080',              // 使わなくても設定しとくと main 側で詰まらない
        // DB 接続（ホスト名とポートは RDS から）
        DB_HOST: db.instanceEndpoint.hostname,
        DB_PORT: db.instanceEndpoint.port.toString(),
        DB_NAME: 'postgres',
        // Secrets の ARN を Lambda に渡す（コード側が ARN を読んで SecretsManager から値を取得）
        DB_SECRET_ARN: dbSecret.secretArn,
        APP_SECRET_ARN: appSecret.secretArn,
        // その他必要な env があればここへ
        CORS_ORIGIN: 'https://word-app-opal.vercel.app',
      },
    });

    // ❶ Secrets 読み取り権限（CDK が GetSecretValue/Describe 用の IAM を付与してくれる）
    dbSecret.grantRead(fn);
    appSecret.grantRead(fn);

    // ❷ RDS への接続（SG で 5432 許可）
    db.connections.allowDefaultPortFrom(fn);

    // VPCE 専用 SG
    // const smVpceSg = new ec2.SecurityGroup(this, 'SmVpceSg', {
    //   vpc: vpc,
    //   description: 'SG for Secrets Manager Interface VPC Endpoint',
    // });

    // VPCE SG に “Lambda SG からの 443” を許可（Inbound は VPCE 側）
    // smVpceSg.addIngressRule(lambdaSg, ec2.Port.tcp(443), 'Lambda to SecretsManager VPCE');

    const shouldCreateSmVpce = !natEnabled && (createSmVpce ?? true);
    if (shouldCreateSmVpce) {
      const smVpceSg = new ec2.SecurityGroup(this, 'SmVpceSg', { vpc, description: 'SG for Secrets Manager VPCE' });
      smVpceSg.addIngressRule(lambdaSg, ec2.Port.tcp(443), 'Lambda to SecretsManager VPCE');
      new ec2.InterfaceVpcEndpoint(this, 'SmEndpoint', {
        vpc,
        service: ec2.InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
        subnets,
        securityGroups: [smVpceSg],
        open: false,              // VPC全許可を無効化
        privateDnsEnabled: true,  // 既定trueだが明示
      });
    }

    // API Gateway（デフォルトは proxy=true で ANY /{proxy+} が生える）
    new apigw.LambdaRestApi(this,'Api',{ handler: fn });
  }
}
