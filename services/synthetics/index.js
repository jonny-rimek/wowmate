var synthetics = require('Synthetics');
const log = require('SyntheticsLogger');

const apiCanaryBlueprint = async function () {

    // Handle validation for positive scenario
    const validateSuccessfull = async function(res) {
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
        hostname: 'bhu3zp80bh.execute-api.us-east-1.amazonaws.com',
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

    await synthetics.executeHttpStep('Verify /combatlogs/keys', requestOptionsStep1, validateSuccessfull, stepConfig1);

    // Set request option for Verify /presign/{filename}
    let requestOptionsStep2 = {
        hostname: 'bhu3zp80bh.execute-api.us-east-1.amazonaws.com',
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

    await synthetics.executeHttpStep('Verify /presign/{filename}', requestOptionsStep2, validateSuccessfull, stepConfig2);

    // Set request option for Verify /combatlogs/keys/{dungeon_id}
    let requestOptionsStep3 = {
        hostname: 'bhu3zp80bh.execute-api.us-east-1.amazonaws.com',
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

    await synthetics.executeHttpStep('Verify /combatlogs/keys/{dungeon_id}', requestOptionsStep3, validateSuccessfull, stepConfig3);

    // Set request option for Verify /combatlogs/keys/{combatlog_uuid}/player-damage-done
    let requestOptionsStep4 = {
        hostname: 'bhu3zp80bh.execute-api.us-east-1.amazonaws.com',
        method: 'GET',
        path: '/combatlogs/keys/fff28fa9-10fb-4018-9486-c1a1f748862d/player-damage-done',
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

    await synthetics.executeHttpStep('Verify /combatlogs/keys/{combatlog_uuid}/player-damage-done', requestOptionsStep4, validateSuccessfull, stepConfig4);


};

exports.handler = async () => {
    return await apiCanaryBlueprint();
};

