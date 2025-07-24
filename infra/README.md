# Welcome to your CDK TypeScript project

This is a blank project for CDK development with TypeScript.

The `cdk.json` file tells the CDK Toolkit how to execute your app.

## Useful commands

* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `npx cdk deploy`  deploy this stack to your default AWS account/region
* `npx cdk diff`    compare deployed stack with current state
* `npx cdk synth`   emits the synthesized CloudFormation template

## 内容の確認(npm run synth後)

#### ファイル全体を less で開く
less cdk.out/DbStack.template.json
####   → 文字列検索
/\"AWS::RDS::DBInstance\"
n   # で次を検索

cat cdk.out/DbStack.template.json | jq '.Resources | keys'
