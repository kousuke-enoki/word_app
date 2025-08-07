import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as ec2 from 'aws-cdk-lib/aws-ec2';

export interface NetworkStackProps extends StackProps {
    natEnabled: boolean;
}
/**
 * ネットワークスタック
 * VPC を定義し、他のスタックから参照できるようにする
 */
export class NetworkStack extends Stack {
  /** 他スタックから参照するため公開 */
  public readonly vpc: ec2.Vpc;

  constructor(scope: Construct, id: string, props: NetworkStackProps) {
    super(scope, id, props);

    this.vpc = new ec2.Vpc(this, 'Vpc', {
      maxAzs: 2,
      natGateways: props.natEnabled ? 1 : 0, // NAT ゲートウェイを有効にするかどうか
      // サブネットの設定
      subnetConfiguration: [
        { name: 'Public', subnetType: ec2.SubnetType.PUBLIC, cidrMask: 24 },
        props.natEnabled
          ? { name: 'AppPrivate', subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS, cidrMask: 24 }
          : { name: 'AppIsolated', subnetType: ec2.SubnetType.PRIVATE_ISOLATED, cidrMask: 24 },
        { name: 'DbIsolated', subnetType: ec2.SubnetType.PRIVATE_ISOLATED, cidrMask: 24 },
      ],
    });
  }
}
