'use strict';

const AWS = require('aws-sdk');
const s3 = new AWS.S3({signatureVersion: 'v4'});
const uuidv4 = require('uuid/v4');

exports.handler = (event, context, callback) => {
	const bucket = process.env['BUCKET_NAME'];
	if (!bucket) {
		callback(new Error(`S3 bucket not set`));
	}

	const path = event.requestContext.http.path;

	let fileending
	let filesize

	if (path.match(/\.txt.gz/) == '.txt.gz') {
		fileending = '.txt.gz';
		//TODO: decrease after figuring out the new limit
		filesize = 100*1024*1024; //100Mibibyte
		//gziped files are smaller, but we still can't process larger files in the lambda
		//once it's uncompressed
		//a 300mb file is around 30mb gziped, but it still has to be unpacked and proccessed,
		//gziping just speeds up the upload and download into the lambda.
		//the size of the file we can handle is limited by the RAM availabe inside the lambda
	} else if (path.match(/\.txt/) == '.txt') {
		fileending = '.txt';
		filesize = 1000*1024*1024; //1000Mibibyte
	} else if (path.match(/\.zip/) == '.zip') {
		fileending = '.zip';
		filesize = 100*1024*1024; //100Mibibyte
	} else {
		callback(null, {
			statusCode: 500,
			// headers: {
			// "Access-Control-Allow-Origin": "*"
			// },
			body: 'invalid filename',
		});
	}

	var currentTime = new Date()

	const key = `${currentTime.getFullYear()}/${currentTime.getMonth() + 1}/${currentTime.getDate()}/${uuidv4()}${fileending}`;

	const res = s3.createPresignedPost({
		Bucket: bucket,
		Fields: {
			key: key,
		},
		Conditions: [
			["content-length-range", 	0, filesize], // content length restrictions: 0-300 MB
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
