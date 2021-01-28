<template>
  <h1>Image Cropper</h1>
  <Cropper @ready="onReady" />
</template>

<script>
import Cropper from "./components/Cropper.vue";

export default {
  components: {
    Cropper,
  },
  methods: {
    onReady(e) {
      console.log("Cropper ready!");
    },
  },
};
</script>

<style></style>
