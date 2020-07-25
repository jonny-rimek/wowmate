<template>
  <div>
    <h1>upload</h1>

    <form
      action="http://localhost:8080/upload"
      method="post"
      enctype="multipart/form-data"
    >
      file: <input type="file" name="file" /><br /><br />
      <input type="submit" value="Submit" />
    </form>
  </div>
</template>

<script>
export default {}
</script>

<style></style>
