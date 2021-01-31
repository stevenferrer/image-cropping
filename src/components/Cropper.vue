<template>
  <div>
    <h3>Input</h3>
    <input type="file" accept="image/*" @change="fileInputChange" />
    <div class="container">
      <div class="image-container">
        <img id="image" class="image" />
      </div>
      <div id="image-preview" class="image-preview image-preview-lg"></div>
    </div>
    <button type="button" @click="cropImage">Crop</button>
    <div class="image-container">
      <h3>Result</h3>
      <img class="image" :src="resultUrl" />
    </div>
  </div>
</template>

<script>
import { onMounted, reactive, ref } from "vue";
import axios from "axios";

import "cropperjs/dist/cropper.css";
import Cropper from "cropperjs";

export default {
  emits: ["ready"],
  setup(_, { emit }) {
    let cropDetail = reactive({});
    let cropper = reactive({});
    let imgFile = reactive({});
    let resultUrl = ref("");

    onMounted(() => {
      const image = document.getElementById("image");
      cropper = new Cropper(image, {
        aspectRatio: 2 / 3,
        viewMode: 3,
        preview: "#image-preview",
        ready(e) {
          emit("ready", e);
        },
        crop(e) {
          cropDetail = e.detail;
        },
      });

      // download a default image
      axios({
        url: "/alicia-vikander.jpg",
        method: "GET",
        responseType: "blob",
      })
        .then(({ data }) => {
          imgFile = data;
          cropper.replace(URL.createObjectURL(imgFile));
        })
        .catch((error) => {
          throw error;
        });
    });

    const cropImage = (e) => {
      if (!(imgFile && cropDetail)) return;

      const formData = new FormData();
      formData.append("image", imgFile);
      formData.append(
        "opts",
        new Blob([JSON.stringify(cropDetail)], { type: "application/json" })
      );
      axios({
        url: "/crop",
        method: "POST",
        data: formData,
        responseType: "blob",
      })
        .then(({ data }) => {
          resultUrl.value = URL.createObjectURL(data);
        })
        .catch((error) => {
          throw error;
        });
    };

    const fileInputChange = (e) => {
      const file = e.target.files[0];
      if (!file || !file.type.includes("image/")) return;

      if (!(typeof FileReader === "function"))
        throw "FileReader API not supported";

      imgFile = file;
      const reader = new FileReader();
      reader.onload = ({ target }) => {
        if (cropper) {
          cropper.replace(target.result);
        }
      };
      reader.readAsDataURL(imgFile);
    };

    return { resultUrl, cropImage, fileInputChange };
  },
};
</script>

<style scoped>
.container {
  display: flex;
  align-items: center;
}

.image-container {
  width: 500px;
  height: 400px;
}

.image {
  display: block;
  max-width: 100%;
}

.image-preview {
  overflow: hidden;
  margin: 0.5em;
}

.image-preview > img {
  max-width: 100%;
}

.image-preview-lg {
  width: 500px;
  height: 400px;
}
</style>
