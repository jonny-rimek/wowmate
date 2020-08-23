<template>
  <div>
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
        thumbnailWidth: 150
      },
      awss3: {
        signingURL: f => {
          return 'https://presign.wowmate.io'
        },
        headers: {
          //   'Content-Type': 'multipart/form-data'
          // doesnt work either
        },
        params: {},
        sendFileToServer: false,
        withCredentials: false
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
