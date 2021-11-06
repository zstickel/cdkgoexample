package main

import (
	"github.com/aws/aws-cdk-go/awscdk"
	"github.com/aws/aws-cdk-go/awscdk/awsappsync"
	"github.com/aws/aws-cdk-go/awscdk/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v3"
	"github.com/aws/jsii-runtime-go"
)

type GocdkStackProps struct {
	awscdk.StackProps
}

func NewGocdkStack(scope constructs.Construct, id string, props *GocdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	// as an example, here's how you would define an AWS SNS topic:

	api := awsappsync.NewGraphqlApi(stack, jsii.String("MyApi"), &awsappsync.GraphqlApiProps{
		AuthorizationConfig: &awsappsync.AuthorizationConfig{
			DefaultAuthorization: &awsappsync.AuthorizationMode{
				AuthorizationType: awsappsync.AuthorizationType_IAM,
			},
		},
		Name:   jsii.String("ZaneApi"),
		Schema: awsappsync.Schema_FromAsset(jsii.String("schema.graphql")),
	})

	table := awsdynamodb.NewTable(stack, jsii.String("Demos"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("id"), Type: awsdynamodb.AttributeType_STRING},
		SortKey:      &awsdynamodb.Attribute{Name: jsii.String("username"), Type: awsdynamodb.AttributeType_STRING},
		Encryption:   awsdynamodb.TableEncryption_DEFAULT,
	})
	demoDS := api.AddDynamoDbDataSource(jsii.String("demoDataSource"), table, &awsappsync.DataSourceOptions{})

	demoDS.CreateResolver(&awsappsync.BaseResolverProps{
		TypeName:                jsii.String("Query"),
		FieldName:               jsii.String("getDemos"),
		RequestMappingTemplate:  awsappsync.MappingTemplate_DynamoDbScanTable(),
		ResponseMappingTemplate: awsappsync.MappingTemplate_DynamoDbResultList(),
	})

	demoDS.CreateResolver(&awsappsync.BaseResolverProps{
		TypeName:                jsii.String("Mutation"),
		FieldName:               jsii.String("addDemo"),
		RequestMappingTemplate:  awsappsync.MappingTemplate_DynamoDbPutItem(awsappsync.PrimaryKey_Partition(jsii.String("id")).Auto(), awsappsync.Values_Projecting(jsii.String("input"))),
		ResponseMappingTemplate: awsappsync.MappingTemplate_DynamoDbResultItem(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewGocdkStack(app, "GocdkStack", &GocdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
