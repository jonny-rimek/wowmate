import cdk = require('@aws-cdk/core');
import * as lambda from '@aws-cdk/aws-lambda';
import apigateway = require('@aws-cdk/aws-apigateway');
import s3 = require('@aws-cdk/aws-s3');
import { HttpApi, LambdaProxyIntegration, HttpMethod } from '@aws-cdk/aws-apigatewayv2';

interface Props extends cdk.StackProps {
	uploadBucket: s3.Bucket;
}

export class Presign extends cdk.Construct {
	public readonly lambda: lambda.Function
	public readonly api: HttpApi;

	constructor(scope: cdk.Construct, id: string, props: Props) {
		super(scope, id)

		const presignLambda = new lambda.Function(this, 'Lambda', {
			runtime: lambda.Runtime.NODEJS_12_X,
			code: lambda.Code.fromAsset('services/upload/presign'),
			handler: 'index.handler',
			environment: {BUCKET_NAME: props.uploadBucket.bucketName},
		});
		this.lambda = presignLambda

		props.uploadBucket.grantPut(presignLambda);

		this.api = new HttpApi(this, 'Api', {
		//TODO: test if i need cors after I activated CORS on the bucket
			corsPreflight: {
				allowOrigins: ["*"],
			},
		})

		this.api.addRoutes({
			path: '/presign',
			methods: [HttpMethod.POST],
			integration: new LambdaProxyIntegration({
				handler: presignLambda,
			}),
		})
	}
}
