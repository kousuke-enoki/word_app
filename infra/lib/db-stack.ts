import { Stack, StackProps, RemovalPolicy, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as rds from 'aws-cdk-lib/aws-rds';
import * as ec2 from 'aws-cdk-lib/aws-ec2';
import * as secrets from 'aws-cdk-lib/aws-secretsmanager';

export interface DbStackProps extends StackProps {
  vpc: ec2.IVpc;
}
/**
 * データベーススタック
 * RDS インスタンスと Secrets Manager のシークレットを定義する
 */
export class DbStack extends Stack {
  public readonly secret: secrets.ISecret;          // ←外部公開
  public readonly db: rds.DatabaseInstance;         // ←外部公開
  public readonly lambdaToDbSecurityGroup: ec2.ISecurityGroup;

  constructor(scope: Construct, id: string, { vpc, ...rest }: DbStackProps) {
    super(scope, id, rest);

    /* ① Secrets */
    this.secret = new rds.DatabaseSecret(this,'Secret',{ username:'postgres' });
    const lambdaToDbSg = new ec2.SecurityGroup(this,'LambdaToDbSG',{
      vpc, description: 'SG for Lambda -> RDS', allowAllOutbound: true,
    });

    /* ② RDS インスタンス */
    this.db = new rds.DatabaseInstance(this,'Rds',{
      vpc,
      vpcSubnets: { subnetType: ec2.SubnetType.PRIVATE_ISOLATED },
      securityGroups: [lambdaToDbSg],
      engine: rds.DatabaseInstanceEngine.postgres({
        version: rds.PostgresEngineVersion.VER_15,
      }),
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T4G, ec2.InstanceSize.MICRO),
      credentials: rds.Credentials.fromSecret(this.secret),
      allocatedStorage: 20,
      removalPolicy: RemovalPolicy.DESTROY,   // 検証用
      backupRetention: Duration.days(1),
    });
    this.db.connections.allowDefaultPortFrom(lambdaToDbSg, 'Lambda access');

    this.lambdaToDbSecurityGroup = lambdaToDbSg;   // ← export
  }
}
