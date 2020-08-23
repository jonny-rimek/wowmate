'use strict';

const AWS = require('aws-sdk');
const s3 = new AWS.S3({signatureVersion: 'v4'});
const uuidv4 = require('uuid/v4');

exports.handler = (event, context, callback) => {
    const bucket = process.env['BUCKET_NAME'];
    if (!bucket) {
        callback(new Error(`S3 bucket not set`));
    }
    console.log(bucket);

    const key = uuidv4() + '.txt';
    const params = {
        'Bucket': bucket,
        'Key': key,
        'Content-Type': '',
    };

	const res = s3.createPresignedPost({
		Bucket: bucket,
		Fields: {
			key: key,
		},
		Conditions: [
			// ["content-length-range", 	0, 1000000000], // content length restrictions: 0-1MB
			//["starts-with", "$Content-Type", "image/"], // content type restriction
		//	["eq", "$x-amz-meta-userid", userid], // tag with userid <= the user can see this!
            {'acl': 'public-read-write'},
            {'success_action_status': '201'},
            ['starts-with', '$Content-Type', ''],
            ['starts-with', '$key', ''],
		]
	});

    let body = {
        signature: {
            'Content-Type': '',
            'acl': 'public-read-write',
            'success_action_status': '201',
            key,
            ...res.fields,
        },
        postEndpoint: res.url,
    }

	callback(null, {
		statusCode: 200,
		headers: {
		  "Access-Control-Allow-Origin": "*"
		},
		body: JSON.stringify(body),
	});
};
