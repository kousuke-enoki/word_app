import { Stack, StackProps } from "aws-cdk-lib";
import { Construct } from "constructs";
import * as ec2 from "aws-cdk-lib/aws-ec2";

export interface NetworkStackProps extends StackProps {
  natEnabled: boolean;
}

export class NetworkStack extends Stack {
  public readonly vpc: ec2.Vpc;

  constructor(scope: Construct, id: string, props: NetworkStackProps) {
    super(scope, id, props);

    this.vpc = new ec2.Vpc(this, "Vpc", {
      maxAzs: 2,
      natGateways: props.natEnabled ? 1 : 0,
      subnetConfiguration: [
        { name: "Public", subnetType: ec2.SubnetType.PUBLIC, cidrMask: 24 },
        //一時的に reserved にして“席”だけ確保しておく
        {
          name: "AppIsolated",
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
          cidrMask: 24,
          reserved: true,
        },
        {
          name: "DbIsolated",
          subnetType: ec2.SubnetType.PRIVATE_ISOLATED,
          cidrMask: 24,
        },
        {
          name: "AppPrivate",
          subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
          cidrMask: 24,
        },
      ],
    });

    const endpointSg = new ec2.SecurityGroup(this, "VpceSg", {
      vpc: this.vpc,
      allowAllOutbound: true,
      description: "VPC endpoints for private AWS APIs",
    });

    // ← ここが重要：NATが無効のときだけ作成、かつコンテキストで切替可能に
    const createAwsApiEndpoints = (this.node.tryGetContext(
      "createAwsApiEndpoints"
    ) ?? !props.natEnabled) as boolean;

    if (createAwsApiEndpoints) {
      // App側のサブネットみに限定（DbIsolated には作らない）
      // const appSubnetGroup = props.natEnabled ? "AppPrivate" : "AppIsolated";
      const appSubnetGroup = "AppPrivate";

      this.vpc.addInterfaceEndpoint("SecretsManagerVPCE", {
        service: ec2.InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
        securityGroups: [endpointSg],
        subnets: { subnetGroupName: appSubnetGroup }, // ここを subnetType から subnetGroupName に
        // privateDnsEnabled: true が既定（重複があると今回のエラー）
      });

      this.vpc.addInterfaceEndpoint("KmsVPCE", {
        service: ec2.InterfaceVpcEndpointAwsService.KMS,
        securityGroups: [endpointSg],
        subnets: { subnetGroupName: appSubnetGroup },
      });
    }
  }
}
