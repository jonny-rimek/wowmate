'use strict';

const AWS = require('aws-sdk');
//TODO: only import the parts we need https://www.youtube.com/watch?v=rrK7PA8ZK7M
//can save atleast 100ms coldstart
const s3 = new AWS.S3({signatureVersion: 'v4'});
const uuidV4 = require('uuid/v4');

exports.handler = (event, context, callback) => {
	const bucket = process.env['BUCKET_NAME'];
	if (!bucket) {
		callback(new Error(`S3 bucket not set`));
	}

	const path = event.requestContext.http.path;

	let fileEnding
	let filesize

	if (path.match(/\.txt.gz/) == '.txt.gz') {
		//ignore IDE warning this needs to be == instead of ===, cus js is fucking garbage
		fileEnding = '.txt.gz';
		//there might be a difference between convert lambda limits and these, because of megabyte vs mebibyte
		filesize = 100*1024*1024; //100Mebibyte
		//gzipped files are smaller, but we still can't process larger files in the lambda
		//once it's uncompressed
		//a 300mb file is around 30mb gzipped, but it still has to be unpacked and processed,
		//zipping just speeds up the upload and download into the lambda.
		//the size of the file we can handle is limited by the RAM available inside the lambda
	} else if (path.match(/\.txt/) == '.txt') {
		fileEnding = '.txt';
		filesize = 1000*1024*1024; //1000Mebibyte
	} else if (path.match(/\.zip/) == '.zip') {
		fileEnding = '.zip';
		filesize = 100*1024*1024; //100Mebibyte
	} else {
		callback(null, {
			statusCode: 500,
			// headers: {
			// "Access-Control-Allow-Origin": "*"
			// },
			body: 'invalid filename',
		});
	}

	const currentTime = new Date()

    //TODO: add upload/ at the front
	const key = `upload/${currentTime.getUTCFullYear()}/${currentTime.getUTCMonth() + 1}/${currentTime.getUTCDay()}/${currentTime.getUTCHours()}/${uuidV4()}${fileEnding}`;

	const res = s3.createPresignedPost({
		Bucket: bucket,
		Fields: {
			key: key,
		},
		Conditions: [
			["content-length-range", 	0, filesize],
			{'success_action_status': '201'},
			['starts-with', '$Content-Type', ''],
			['starts-with', '$key', ''],
		]
	});

	let body = {
		signature: {
			'Content-Type': '',
			'success_action_status': '201',
			key,
			...res.fields,
		},
		postEndpoint: res.url,
	}

	callback(null, {
		statusCode: 200,
		// headers: {
		//   "Access-Control-Allow-Origin": "*"
		// },
		body: JSON.stringify(body),
	});
};
