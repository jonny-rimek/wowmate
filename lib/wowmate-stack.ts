// import { Frontend } from './frontend';
import { V2 } from './v2';
import { Construct, Stage, Stack, StackProps, StageProps, SecretValue } from '@aws-cdk/core';
import { CdkPipeline, SimpleSynthAction } from '@aws-cdk/pipelines';
import * as codepipeline from '@aws-cdk/aws-codepipeline';
import * as codepipeline_actions from '@aws-cdk/aws-codepipeline-actions'

export class Wowmate extends Stage {
	constructor(scope: Construct, id: string, props?: StageProps) {
		super(scope, id, props);

		new V2(this, 'V2')
		// new Frontend(this, 'frontend')
	}
}

export class WowmatePipelineStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const sourceArtifact = new codepipeline.Artifact();
    const cloudAssemblyArtifact = new codepipeline.Artifact();

    const pipeline = new CdkPipeline(this, 'Pipeline', {
      pipelineName: 'WowmatePipeline',
      cloudAssemblyArtifact,

      sourceAction: new codepipeline_actions.GitHubSourceAction({
		actionName: 'GitHub',
		output: sourceArtifact,
		branch: 'master',
		oauthToken: SecretValue.secretsManager('github-personal-access-token'),
		//TODO: switch to webhook, might need to update the oath token
        trigger: codepipeline_actions.GitHubTrigger.POLL,
        owner: 'jonny-rimek',
        repo: 'wowmate',
      }),
      synthAction: SimpleSynthAction.standardNpmSynth({
        sourceArtifact,
        cloudAssemblyArtifact,

        buildCommand: 'npm run build',
      }),
    });

	const wmp = pipeline.addApplicationStage(new Wowmate(this, 'Prod'));
	wmp.addManualApprovalAction();
  }
}
/* 
export class Wowmate extends cdk.Stack {
	constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
		super(scope, id, props);

		// const db = new Database(this, 'Database')

		// new Api(this, 'Api', {
		// 	dynamoDB: db.dynamoDB,
		// })

		// new Frontend(this, 'Frontend')
		
		// new Upload(this, 'Upload', {
		// 	dynamoDB: db.dynamoDB,
		// })

		new V2(this, 'V2')
		new Frontend(this, 'frontend')
	}
}
 */