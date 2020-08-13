#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { V2 } from '../lib/v2';
import { Construct, Stage, Stack, StackProps, StageProps, SecretValue } from '@aws-cdk/core';
import { CdkPipeline, SimpleSynthAction } from '@aws-cdk/pipelines';
import * as codepipeline from '@aws-cdk/aws-codepipeline';
import * as codepipeline_actions from '@aws-cdk/aws-codepipeline-actions'

class Wowmate extends Stage {
	constructor(scope: Construct, id: string, props?: StageProps) {
		super(scope, id, props);

		new V2(this, 'V2')
		// new Frontend(this, 'frontend')
	}
}

class WowmatePipelineStack extends Stack {
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
TODO: add 
interface EnvProps {
  prod: boolean;
}
update constructors 
constructor(scope: Construct, id: string, props?: EnvProps) {
	*/

const app = new cdk.App();
new WowmatePipelineStack(app, 'WoWM', {
	env: {region: "us-east-1"}
});