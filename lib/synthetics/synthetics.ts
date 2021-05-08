import * as synthetics from '@aws-cdk/aws-synthetics';
import { SnsAction } from '@aws-cdk/aws-cloudwatch-actions';
import * as kms from "@aws-cdk/aws-kms";
import sns = require('@aws-cdk/aws-sns');
import * as subscriptions from '@aws-cdk/aws-sns-subscriptions';
import * as iam from '@aws-cdk/aws-iam';
import cdk = require('@aws-cdk/core');
import cloudwatch = require('@aws-cdk/aws-cloudwatch');
import s3 = require('@aws-cdk/aws-s3');

interface Props extends cdk.StackProps {
    errorMail: string
    stage: string
    apiUrl: string
    domainName: string
    key: kms.Key
}

export class Synthetics extends cdk.Construct {

    constructor(scope: cdk.Construct, id: string, props: Props) {
        super(scope, id)
        const code = `
            var synthetics = require('Synthetics');
            const log = require('SyntheticsLogger');

            const apiCanaryBlueprint = async function () {
                const hostname = process.env.API_URL
                const domainName = process.env.DOMAIN_NAME
                
                // Handle validation for positive scenario
                const validateSuccessful = async function(res) {
                    return new Promise((resolve, reject) => {
                        if (res.statusCode < 200 || res.statusCode > 299) {
                            throw res.statusCode + ' ' + res.statusMessage;
                        }

                        let responseBody = '';
                        res.on('data', (d) => {
                            responseBody += d;
                        });

                        res.on('end', () => {
                            // Add validation on 'responseBody' here if required.
                            resolve();
                        });
                    });
                };


                // Set request option for Verify /combatlogs/keys
                let requestOptionsStep1 = {
                    hostname: hostname,
                    method: 'GET',
                    path: '/combatlogs/keys',
                    port: '443',
                    protocol: 'https:',
                    body: "",
                    headers: {}
                };
                requestOptionsStep1['headers']['User-Agent'] = [synthetics.getCanaryUserAgentString(), requestOptionsStep1['headers']['User-Agent']].join(' ');

                // Set step config option for Verify /combatlogs/keys
                let stepConfig1 = {
                    includeRequestHeaders: false,
                    includeResponseHeaders: false,
                    includeRequestBody: false,
                    includeResponseBody: false,
                    restrictedHeaders: [],
                    continueOnHttpStepFailure: true
                };

                await synthetics.executeHttpStep('Verify /combatlogs/keys', requestOptionsStep1, validateSuccessful, stepConfig1);

                // Set request option for Verify /presign/{filename}
                let requestOptionsStep2 = {
                    hostname: hostname,
                    method: 'POST',
                    path: '/presign/test.txt',
                    port: '443',
                    protocol: 'https:',
                    body: "",
                    headers: {}
                };
                requestOptionsStep2['headers']['User-Agent'] = [synthetics.getCanaryUserAgentString(), requestOptionsStep2['headers']['User-Agent']].join(' ');

                // Set step config option for Verify /presign/{filename}
                let stepConfig2 = {
                    includeRequestHeaders: false,
                    includeResponseHeaders: false,
                    includeRequestBody: false,
                    includeResponseBody: false,
                    restrictedHeaders: [],
                    continueOnHttpStepFailure: true
                };

                await synthetics.executeHttpStep('Verify /presign/{filename}', requestOptionsStep2, validateSuccessful, stepConfig2);

                // Set request option for Verify /combatlogs/keys/{dungeon_id}
                let requestOptionsStep3 = {
                    hostname: hostname,
                    method: 'GET',
                    path: '/combatlogs/keys/2291',
                    port: '443',
                    protocol: 'https:',
                    body: "",
                    headers: {}
                };
                requestOptionsStep3['headers']['User-Agent'] = [synthetics.getCanaryUserAgentString(), requestOptionsStep3['headers']['User-Agent']].join(' ');

                // Set step config option for Verify /combatlogs/keys/{dungeon_id}
                let stepConfig3 = {
                    includeRequestHeaders: false,
                    includeResponseHeaders: false,
                    includeRequestBody: false, 
                    includeResponseBody: false,
                    restrictedHeaders: [],
                    continueOnHttpStepFailure: true
                };

                // TODO: XMLHttpRequest
                // TODO: check response, must include body
                await synthetics.executeHttpStep('Verify /combatlogs/keys/{dungeon_id}', requestOptionsStep3, validateSuccessful, stepConfig3);
                
                const combatlogUUID = getCombatlogUUID(hostname)
                // Set request option for Verify /combatlogs/keys/{combatlog_uuid}/player-damage-done
                let requestOptionsStep4 = {
                    hostname: hostname,
                    method: 'GET',
                    path: '/combatlogs/keys/' + combatlogUUID,
                    port: '443',
                    protocol: 'https:',
                    body: "",
                    headers: {}
                };
                requestOptionsStep4['headers']['User-Agent'] = [synthetics.getCanaryUserAgentString(), requestOptionsStep4['headers']['User-Agent']].join(' ');

                // Set step config option for Verify /combatlogs/keys/{combatlog_uuid}/player-damage-done
                let stepConfig4 = {
                    includeRequestHeaders: false,
                    includeResponseHeaders: false,
                    includeRequestBody: false,
                    includeResponseBody: false,
                    restrictedHeaders: [],
                    continueOnHttpStepFailure: true
                };

                await synthetics.executeHttpStep('Verify /combatlogs/keys/{combatlog_uuid}/player-damage-done', requestOptionsStep4, validateSuccessful, stepConfig4);
                
                // ----------------------------------------------------
                // Set request option for Verify /
                let requestOptionsStep5 = {
                    hostname: domainName,
                    method: 'GET',
                    path: '/',
                    port: '443',
                    protocol: 'https:',
                    body: "",
                    headers: {}
                };
                requestOptionsStep4['headers']['User-Agent'] = [synthetics.getCanaryUserAgentString(), requestOptionsStep4['headers']['User-Agent']].join(' ');

                // Set step config option for Verify /
                let stepConfig5 = {
                    includeRequestHeaders: false,
                    includeResponseHeaders: false,
                    includeRequestBody: false,
                    includeResponseBody: false,
                    restrictedHeaders: [],
                    continueOnHttpStepFailure: true
                };

                await synthetics.executeHttpStep('Verify /', requestOptionsStep5, validateSuccessful, stepConfig5);
            };
            
            // neither fetch or XMLHttpRequest works and I can't install packages with inline code
            function getCombatlogUUID(apiUrl) {
                if (process.env.STAGE === "prod") {
                    return "ac18f293-e573-4994-8745-2a291d12bf9a"
                }
                return "2ac44ef5-d053-4232-881e-c8c21300ec4f"
                // let url = "https://" + apiUrl + "/combatlogs/keys"
                // let xmlHttp = new XMLHttpRequest()
                // xmlHttp.open( "GET", url, false ) // false for synchronous request
                // xmlHttp.send( null )
                // const resp = JSON.parse(xmlHttp.responseText)
                // return resp.data[0].combatlog_uuid
            }

            exports.handler = async () => {
                return await apiCanaryBlueprint();
            };
        `

        let schedule = synthetics.Schedule.rate(cdk.Duration.minutes(5))
        if (props.stage == "preprod" || props.stage == "dev") {
            // u can't invoke the canary manually if it runs on a schedule and I don't need the canary on a schedule in dev and preprod anyway
            schedule = synthetics.Schedule.once()
        }

        const canary = new synthetics.Canary(this, 'API', {
            schedule: schedule,
            runtime: synthetics.Runtime.SYNTHETICS_NODEJS_PUPPETEER_3_0,
            test: synthetics.Test.custom({
                handler: 'index.handler',
                code: synthetics.Code.fromInline(code),
            }),
        });
        const cfnCanary = canary.node.defaultChild as synthetics.CfnCanary
        cfnCanary.addPropertyOverride('RunConfig.EnvironmentVariables', {
            API_URL: props.apiUrl.replace('https://','').replace('/',''),
            STAGE: props.stage,
            DOMAIN_NAME: props.domainName,
        })

        const cfnBucket = canary.artifactsBucket.node.defaultChild as s3.CfnBucket
        cfnBucket.cfnOptions.metadata = {
            cfn_nag: {
                rules_to_suppress: [
                    {
                        id: 'W35',
                        reason: "this is a canary bucket, so there is no point tracking access to it",
                    },
                    {
                        id: 'W51',
                        reason: "this is a canary bucket, it doesn't need a bucket policy",
                    },
                ]
            }
        }
        const cfnRole = canary.role.node.defaultChild as iam.CfnRole
        cfnRole.cfnOptions.metadata = {
            cfn_nag: {
                rules_to_suppress: [
                    {
                        id: 'W11',
                        reason: "this is a CDK role #trust",
                    },
                ]
            }
        }

        const errorTopic = new sns.Topic(this, 'errorTopic', {
        	masterKey: props.key
        });
        errorTopic.addSubscription(new subscriptions.EmailSubscription(props.errorMail));

        new cloudwatch.Alarm(this, 'CanaryAlarm', {
            metric: canary.metricSuccessPercent(),
            evaluationPeriods: 2,
            threshold: 100,
            comparisonOperator: cloudwatch.ComparisonOperator.LESS_THAN_THRESHOLD,
        }).addAlarmAction(new SnsAction(errorTopic))


    }
}
