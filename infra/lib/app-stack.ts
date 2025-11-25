import {
  Duration,
  RemovalPolicy,
  Stack,
  StackProps,
  CfnOutput,
} from "aws-cdk-lib";
import * as apigw from "aws-cdk-lib/aws-apigateway";
import * as cloudwatch from "aws-cdk-lib/aws-cloudwatch";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as ecr from "aws-cdk-lib/aws-ecr";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as logs from "aws-cdk-lib/aws-logs";
import * as rds from "aws-cdk-lib/aws-rds";
import * as secrets from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";

export interface AppStackProps extends StackProps {
  vpc: ec2.IVpc;
  secret: secrets.ISecret;
  // DB 用 Secret
  db: rds.IDatabaseInstance;
  lambdaSg: ec2.ISecurityGroup;
  natEnabled: boolean;
  appSecretArn?: string;
  // APP 用 Secret の ARN を受け取れるようにする
  createSmVpce?: boolean;
}

/**
 * アプリケーションスタック
 * ECR イメージを Lambda で実行し、API Gateway で公開する
 */
export class AppStack extends Stack {
  constructor(
    scope: Construct,
    id: string,
    {
      vpc,
      secret: dbSecret,
      db,
      lambdaSg,
      natEnabled,
      appSecretArn,
      createSmVpce,
      ...rest
    }: AppStackProps
  ) {
    super(scope, id, rest);

    const appSecret = secrets.Secret.fromSecretCompleteArn(
      this,
      "AppSecretImported",
      appSecretArn ??
        "arn:aws:secretsmanager:ap-northeast-1:381492105871:secret:wordapp/app-k4I6ng"
    );

    // VPC 内のサブネット（NAT Gateway 有無で変える）
    // const subnets = natEnabled
    //   ? { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS }
    //   : { subnetType: ec2.SubnetType.PRIVATE_ISOLATED };
    const subnets = { subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS };

    const imageTagOrDigest =
      this.node.tryGetContext("imageDigest") ??
      this.node.tryGetContext("imageTag") ??
      "lambda";

    // ECR イメージ → Lambda
    const repo = ecr.Repository.fromRepositoryName(
      this,
      "Repo",
      "wordapp-backend"
    );
    const fn = new lambda.DockerImageFunction(this, "ApiFn", {
      code: lambda.DockerImageCode.fromEcr(repo, {
        tagOrDigest: imageTagOrDigest,
      }),
      vpc,
      vpcSubnets: { subnetGroupName: "AppPrivate" },
      securityGroups: [lambdaSg],
      memorySize: 256,
      timeout: Duration.seconds(30), // ← 少し余裕を持たせる
      // reservedConcurrentExecutions: 10, // 予約同時実行数
      // ↑ アカウント全体の “未予約(Unreserved) 同時実行数” を最低必要値 10 未満に落としてしまったため一旦削除
      logRetention: logs.RetentionDays.THREE_DAYS,
      environment: {
        APP_ENV: "production",
        GIN_MODE: "release",
        APP_PORT: "8080", // 使わなくても設定しとくと main 側で詰まらない
        // DB 接続（ホスト名とポートは RDS から）
        DB_HOST: db.instanceEndpoint.hostname,
        DB_PORT: db.instanceEndpoint.port.toString(),
        DB_NAME: "postgres",
        DB_USER: dbSecret.secretValueFromJson("username").unsafeUnwrap(),
        DB_PASSWORD: dbSecret.secretValueFromJson("password").unsafeUnwrap(),
        LINE_CLIENT_ID: appSecret
          .secretValueFromJson("LINE_CLIENT_ID")
          .unsafeUnwrap(),
        LINE_CLIENT_SECRET: appSecret
          .secretValueFromJson("LINE_CLIENT_SECRET")
          .unsafeUnwrap(),
        LINE_REDIRECT_URI: appSecret
          .secretValueFromJson("LINE_REDIRECT_URI")
          .unsafeUnwrap(),
        JWT_SECRET: appSecret.secretValueFromJson("JWT_SECRET").unsafeUnwrap(),
        // Secrets の ARN を Lambda に渡す（コード側が ARN を読んで SecretsManager から値を取得）
        DB_SECRET_ARN: dbSecret.secretArn,
        APP_SECRET_ARN: appSecret.secretArn,
        CORS_ORIGIN: "https://word-app-opal.vercel.app",
        // 起動時の重さ回避
        // 注意: 初回デプロイ時やスキーマ変更時は "true" に設定してMigrationを実行
        // Migration完了後は "false" に戻すことを推奨
        RUN_MIGRATION: "true",
        RUN_SEEDER: "true",
        RUN_SEEDER_FOR_WORDS: "false",
        APP_BOOTSTRAP_MODE: "FULL",
        ENABLE_TEST_USER_MODE: "true",
        // ↓再デプロイの「差分」が毎回出るので、
        // secret managerの変更後確実に再構成される
        DEPLOY_REV: Date.now().toString(),
      },
    });

    // ❶ Secrets 読み取り権限（CDK が GetSecretValue/Describe 用の IAM を付与してくれる）
    // dbSecret.grantRead(fn);
    // appSecret.grantRead(fn);
    // ❷ RDS への接続（SG で 5432 許可）
    db.connections.allowDefaultPortFrom(fn);
    // VPCE SG に “Lambda SG からの 443” を許可（Inbound は VPCE 側）
    // smVpceSg.addIngressRule(lambdaSg, ec2.Port.tcp(443),'Lambda to SecretsManager VPCE');
    const shouldCreateSmVpce = !natEnabled && (createSmVpce ?? true);

    if (shouldCreateSmVpce) {
      const smVpceSg = new ec2.SecurityGroup(this, "SmVpceSg", {
        vpc,
        description: "SG for Secrets Manager VPCE",
      });
      smVpceSg.addIngressRule(
        lambdaSg,
        ec2.Port.tcp(443),
        "Lambda to SecretsManager VPCE"
      );
      new ec2.InterfaceVpcEndpoint(this, "SmEndpoint", {
        vpc,
        service: ec2.InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
        subnets,
        securityGroups: [smVpceSg],
        open: false, // VPC全許可を無効化
        privateDnsEnabled: true, // 既定trueだが明示 });
      });
    }

    // DynamoDB テーブル（レート制限用）
    const rateLimitsTable = new dynamodb.Table(this, "RateLimitsTable", {
      tableName: "rate_limits",
      partitionKey: { name: "pk", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "sk", type: dynamodb.AttributeType.STRING },
      timeToLiveAttribute: "ttl",
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: RemovalPolicy.DESTROY, // 本番は RETAIN 推奨
    });

    // Lambda に DynamoDB ReadWrite 権限を付与
    rateLimitsTable.grantReadWriteData(fn);

    fn.addEnvironment("RATE_LIMIT_TABLE", rateLimitsTable.tableName);
    // レート制限バックエンド（DynamoDBを使用）
    fn.addEnvironment("RATE_LIMIT_BACKEND", "dynamodb");
    // レート制限設定
    fn.addEnvironment("RATE_LIMIT_WINDOW_SEC", "60");
    fn.addEnvironment("RATE_LIMIT_MAX_REQUESTS", "1");
    fn.addEnvironment("RATE_LIMIT_TTL_SEC", "120");
    // リミット設定
    fn.addEnvironment("LIMIT_REGISTERED_WORDS_PER_USER", "200");
    fn.addEnvironment("LIMIT_QUIZ_MAX_PER_DAY", "20");
    fn.addEnvironment("LIMIT_QUIZ_MAX_QUESTIONS", "100");
    fn.addEnvironment("LIMIT_BULK_MAX_PER_DAY", "5");
    fn.addEnvironment("LIMIT_BULK_MAX_BYTES", "51200");
    fn.addEnvironment("LIMIT_BULK_TOKENIZE_MAX_TOKENS", "51200");
    fn.addEnvironment("LIMIT_BULK_REGISTER_MAX_ITEMS", "51200");

    // API Gateway の作成（スロットリング設定付き）
    // 注意: CORSはGin側（バックエンド）で処理するため、API Gateway側では設定しない
    const api = new apigw.RestApi(this, "Api", {
      restApiName: `${this.stackName}-api`,
      description: "API for word app",
      deployOptions: {
        stageName: "prod",
        throttlingRateLimit: 50,
        throttlingBurstLimit: 100,
        // メソッド別スロットリングを追加
        methodOptions: {
          // キーはリソースパス/HTTPメソッド
          "/auth/test-login/POST": {
            throttlingRateLimit: 5,
            throttlingBurstLimit: 10,
          },
        },
      },
    });
    // ルート統合（プロキシ統合）
    const proxyResource = api.root.addResource("{proxy+}");
    proxyResource.addMethod(
      "ANY",
      new apigw.LambdaIntegration(fn, {
        proxy: true,
      })
    );

    // /auth/test-login エンドポイント
    const testLoginResource = api.root
      .addResource("users")
      .addResource("auth")
      .addResource("test-login");
    testLoginResource.addMethod(
      "POST",
      new apigw.LambdaIntegration(fn, {
        proxy: true,
      }),
      {
        methodResponses: [
          {
            statusCode: "200",
          },
          {
            statusCode: "429",
          },
          {
            statusCode: "500",
          },
        ],
      }
    );

    // Usage Plan で /auth/test-login エンドポイントに個別のスロットリング設定
    // 注意: API Gateway REST API では Stage レベルでのスロットリングしか設定できないため、
    // 個別エンドポイントの制限はアプリ側（DynamoDB）で実装
    // ここでは Stage 全体のスロットリング（50/s, burst=100）のみ設定
    // より細かい制御が必要な場合は Usage Plan + API Key を使用

    // CloudWatch Alarm（5xx エラー監視）
    // API Gateway の 5xx エラーメトリクスを監視
    const errorAlarm = new cloudwatch.Alarm(this, "Api5xxAlarm", {
      alarmName: `${this.stackName}-api-5xx-errors`,
      metric: api.metricServerError({
        period: Duration.minutes(5),
        statistic: "Sum",
      }),
      threshold: 50,
      evaluationPeriods: 1,
      treatMissingData: cloudwatch.TreatMissingData.NOT_BREACHING,
      alarmDescription:
        "Alarm when API Gateway 5xx errors exceed 50 in 5 minutes",
    });

    // アラーム通知（任意: SNSトピックが必要な場合は作成）
    // const alarmTopic = new sns.Topic(this, "AlarmTopic");
    // errorAlarm.addAlarmAction(new cloudwatch_actions.SnsAction(alarmTopic));

    // API Gateway URL を Output として出力
    new CfnOutput(this, "ApiUrl", {
      value: api.url,
      description: "API Gateway URL",
      exportName: `${this.stackName}-ApiUrl`,
    });
  }
}
