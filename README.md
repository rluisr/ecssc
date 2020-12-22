ecssc
======
ecssc(ECS State Check) is a Lambda function for notification to Slack if the ECS task event is changed.

![image](https://f.easyuploader.app/eu-prd/upload/20201213005728_6f5264366f3754797957.png)

Download
---------
See [GitHub Container Registry](https://github.com/users/rluisr/packages/container/ecssc/versions)

Installation
-------------
ecssc is a Lambda function of EventBridge event target.

![image](https://f.easyuploader.app/eu-prd/upload/20201213011250_374537644f5a66646969.png)

**Copy the container image to your ECR repository.**  
The Lambda does not support third-party container registry.

This is an example for setting up EventBridge with Terraform.

EventBridge:
```hcl
resource "aws_cloudwatch_event_bus" "ecs-state-check" {
  name = "ecs-state-check"
}

resource "aws_cloudwatch_event_rule" "ecs-state-check" {
  name        = "ecs-state-check"
  event_pattern = <<EOF
{
  "source": [
    "aws.ecs"
  ],
  "detail-type": [
    "ECS Task State Change",
    "ECS Container Instance State Change"
  ]
}
EOF
}
```

Lambda:
```hcl
# ECR
resource "aws_ecr_repository" "foo" {
  name                 = "ecs-state-check"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = false
  }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "ecs-state-check" {
  name              = "/aws/lambda/ecs-state-check"
  retention_in_days = 14
}

# IAM role for Lambda
module "ecs-state-check_lambda_execution_role" {
  source  = "baikonur-oss/iam-nofile/aws"
  version = "v2.0.0"

  type = "lambda"
  name = "ecs-state-check"

  policy_json = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ssm:GetParametersByPath"
      ],
      "Resource": "arn:aws:ssm:${data.aws_region.region.name}:${data.aws_caller_identity.caller.account_id}:parameter${var.parameter_store_path}"
    },
    {
      "Effect": "Allow",
      "Action": [
        "ecs:Describe*"
      ],
      "Resource": "*"
    }
  ]
}
EOF
}

# Lambda
resource "aws_lambda_function" "ecs-state-check" {
  function_name = "ecs-state-check"
  image_uri = "to your ecr image uri"
  package_type = "Image"
  role = module.ecs-state-check_lambda_execution_role.arn
}

# Trigger
resource "aws_lambda_permission" "ecs-state-check" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.ecs-state-check.function_name
  principal     = "events.amazonaws.com"
  source_arn    = "arn:aws:events:ap-northeast-1:111111111111:rule/ecs-state-check"
}
```

Environment Variables
----------------------
| Name                         | Description                                        | Required |
|------------------------------|----------------------------------------------------|----------|
| ECSSC_DEBUG                  | Show debug message. default: false                 | no       |
| ECSSC_IGNORE_CONTAINER_NAMES | Skip container names. Support array like app1,app2 | no       |
| ECSSC_SLACK_CHANNEL_NAME     | Slack channel name like '#test'                    | yes      |
| ECSSC_SLACK_WEBHOOK_URL      | Slack incoming webhook URL                         | yes      |
| ECSSC_SLACK_USER_NAME        | Slack Username default 'ecs-state-check'           | no       |
| ECSSC_SLACK_ICON_URL         | Slack icon URL                                     | no       |
| ECSSC_SLACK_ICON_EMOJI       | Slack icon emoji. default ':japanese_goblin:'      | no       |