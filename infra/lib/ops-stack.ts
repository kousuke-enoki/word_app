import { Stack, StackProps, Duration } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as budgets from 'aws-cdk-lib/aws-budgets';
import * as sns     from 'aws-cdk-lib/aws-sns';
import * as subs    from 'aws-cdk-lib/aws-sns-subscriptions';
import * as lambda  from 'aws-cdk-lib/aws-lambda';
import * as events  from 'aws-cdk-lib/aws-events';
import * as targets from 'aws-cdk-lib/aws-events-targets';
import * as rds     from 'aws-cdk-lib/aws-rds';
import * as iam from 'aws-cdk-lib/aws-iam';

export interface OpsStackProps extends StackProps {
  db: rds.DatabaseInstance;              // DbStack から渡す
  stopHourJst?: number;                  // 例: 0 (= 深夜 0 時)
  startHourJst?: number;                 // 例: 9 (= 朝 9 時)
}

export class OpsStack extends Stack {
  constructor(scope: Construct, id: string, props: OpsStackProps) {
    super(scope, id, props);


    /* デフォルト時刻 */
    const stopAt  = props.stopHourJst  ?? 0;  // JST 00:00
    const startAt = props.startHourJst ?? 9;  // JST 09:00

    /* ① Lambda – Start / Stop ロジック（100 行も要らない簡素版） */
    const scheduler = new lambda.Function(this, 'RdsSchedulerFn', {
      runtime : lambda.Runtime.PYTHON_3_12,
      handler : 'index.handler',
      timeout : Duration.seconds(60),
      memorySize: 128,
      description: 'Stop RDS at night, start in the morning (JST)',
      code: lambda.Code.fromInline(`
        import boto3, os, datetime
        rds = boto3.client('rds')
        DB  = os.environ['DB_ID']
        STOP = int(os.environ['STOP_HOUR'])
        START = int(os.environ['START_HOUR'])

        def handler(event, _):
            # JST = UTC+9
            jst_hour = (datetime.datetime.utcnow().hour + 9) % 24

            if STOP <= jst_hour < START:
                print(f"Stopping {DB}")
                rds.stop_db_instance(DBInstanceIdentifier=DB)
            elif jst_hour == START:
                print(f"Starting {DB}")
                rds.start_db_instance(DBInstanceIdentifier=DB)
            else:
                print(f"No action needed at {jst_hour} JST")
        `),
      environment: {
        DB_ID      : props.db.instanceIdentifier,
        STOP_HOUR  : stopAt.toString(),
        START_HOUR : startAt.toString(),
      },
      logRetention: 3,  // 日数 / aws-logs RetentionDays enum でも可
    });

    /* ② 必要 IAM 権限を Lambda に付与 */
    scheduler.addToRolePolicy(new iam.PolicyStatement({
      actions: ['rds:StartDBInstance', 'rds:StopDBInstance'],
      resources: [props.db.instanceArn],
    }));

    /* ③ EventBridge ルール ― 1 時間おきに起動 */
    new events.Rule(this, 'RunEveryHour', {
      schedule: events.Schedule.rate(Duration.hours(1)),
      targets : [new targets.LambdaFunction(scheduler)],
    });

    /* 1. 料金アラート */
    const topic = new sns.Topic(this, 'BillingTopic');
    topic.addSubscription(new subs.EmailSubscription('billing@example.com'));

    new budgets.CfnBudget(this, 'CostBudget', {
      budget: {
        budgetType : 'COST',
        timeUnit   : 'MONTHLY',
        budgetLimit: { amount: 1000, unit: 'JPY' },
      },
      notificationsWithSubscribers: [{
        notification: {
          notificationType  : 'ACTUAL',
          comparisonOperator: 'GREATER_THAN',
          threshold         : 100,
          thresholdType     : 'PERCENTAGE',
        },
        subscribers: [{
          subscriptionType: 'SNS',
          address        : topic.topicArn,
        }],
      }],
    });

    /* 2. RDS 自動停止／起動 Lambda */
    const schedFn = new lambda.Function(this, 'RdsScheduler', {
      runtime : lambda.Runtime.PYTHON_3_12,
      code    : lambda.Code.fromInline(`
        import boto3, os, datetime, pytz
        rds = boto3.client('rds')
        DB = os.environ['DB_ID']
        TZ = pytz.timezone('Asia/Tokyo')
        def handler(e,_):
            h = datetime.datetime.now(TZ).hour
            if 0 <= h < 8:
                rds.stop_db_instance(DBInstanceIdentifier=DB)
            elif h == 9:
                rds.start_db_instance(DBInstanceIdentifier=DB)
        `),
      handler : 'index.handler',
      timeout : Duration.seconds(60),
      environment: { DB_ID: props.db.instanceIdentifier },
    });
    schedFn.addToRolePolicy(new iam.PolicyStatement({
      actions: ['rds:StartDBInstance', 'rds:StopDBInstance'],
      resources: [props.db.instanceArn],      // その DB インスタンスのみ
    }));
    // db.grantStopStart(schedFn);

    new events.Rule(this, 'Hourly', {
      schedule: events.Schedule.rate(Duration.hours(1)),
      targets : [new targets.LambdaFunction(schedFn)],
    });
  }
}
