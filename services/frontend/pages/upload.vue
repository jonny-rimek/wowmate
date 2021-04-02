<template>
	<div>
		<div
			class="md:flex md:items-center md:justify-between pt-4 lg:pt-8 pb-4 lg:pb-6 pl-4 lg:pl-0"
		>
			<div class="flex-1 min-w-0">
				<h2
					class="text-l font-bold leading-7 text-gray-200 sm:text-3xl sm:leading-9 sm:truncate"
				>
					Combatlog Upload
				</h2>
			</div>
		</div>

		<Dropzone
			id="dropzone"
			:awss3="awss3"
			:options="dropzoneOptions"
			@vdropzone-s3-upload-error="s3UploadError"
			@vdropzone-s3-upload-success="s3UploadSuccess"
			@vdropzone-success="dropzoneSuccess"
		/>

		Note: doesn't work with files larger than 100MB
	</div>
</template>

<script>
import Dropzone from 'nuxt-dropzone'
import 'nuxt-dropzone/dropzone.css'

export default {
	components: {
		Dropzone
	},
	data() {
		return {
			url: 'default',
			images: {},
			dropzoneOptions: {
				method: 'POST',
				thumbnailWidth: 150,
				maxFileSize: 1000
			},
			awss3: {
				signingURL: f => {
					return process.env.baseUrl + '/presign/' + f.name
				},
				headers: {},
				params: {},
				sendFileToServer: false,
				withCredentials: false,

			}
		}
	},
	methods: {
		dropzoneSuccess(file, res) {
			console.log('success?')
		},
		s3UploadError(errorMessage) {
			// Show an error message on failure
			console.log(errorMessage)
		},
		s3UploadSuccess(s3ObjectLocation) {
			// Show a message after uploaded to S3?
			console.log(s3ObjectLocation)
		}
	}
}
</script>

<style></style>
