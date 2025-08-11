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
  public readonly vpc: ec2.Vpc;

  constructor(scope: Construct, id: string, props: NetworkStackProps) {
    super(scope, id, props);

    this.vpc = new ec2.Vpc(this, 'Vpc', {
      maxAzs: 2,
      natGateways: props.natEnabled ? 1 : 0, // NAT ゲートウェイを有効にするかどうか
      subnetConfiguration: [
        { name: 'Public', subnetType: ec2.SubnetType.PUBLIC, cidrMask: 24 },
        props.natEnabled
          ? { name: 'AppPrivate', subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS, cidrMask: 24 }
          : { name: 'AppIsolated', subnetType: ec2.SubnetType.PRIVATE_ISOLATED, cidrMask: 24 },
        { name: 'DbIsolated', subnetType: ec2.SubnetType.PRIVATE_ISOLATED, cidrMask: 24 },
      ],
    });

    // NATなしで AWS API(Secrets Manager) を叩くための VPC Endpoint
    const vpceSg = new ec2.SecurityGroup(this, 'VpceSg', {
      vpc: this.vpc,
      description: 'VPC endpoints for private AWS APIs',
      allowAllOutbound: true,
    });
    this.vpc.addInterfaceEndpoint('SecretsManagerEndpoint', {
      service: ec2.InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
      subnets: {
        subnetGroupName: props.natEnabled ? 'AppPrivate' : 'AppIsolated',
      },
      securityGroups: [vpceSg],
    });
    // ついでに SSM / STS も使うなら追加（任意）
    // this.vpc.addInterfaceEndpoint('SsmEndpoint', { service: ec2.InterfaceVpcEndpointAwsService.SSM, subnets:{ subnetGroupName: ... }, securityGroups:[vpceSg] });
    // this.vpc.addInterfaceEndpoint('StsEndpoint', { service: ec2.InterfaceVpcEndpointAwsService.STS, subnets:{ subnetGroupName: ... }, securityGroups:[vpceSg] });
  }
}

