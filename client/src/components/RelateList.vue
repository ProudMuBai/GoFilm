<template>
  <div class="container">
      <h2 class="title">相关影片</h2>
      <div class="content">
          <div class="item" v-for="m in relateList">
              <a :href="`/filmDetail?link=${m.id}`" class="cus_content_link"
                 :style="{backgroundImage: `url('${m.picture}')`}">
                  <span class="cus_tag hidden-md-and-down">{{ m.year }}</span>
                  <span class="cus_tag hidden-md-and-down">{{ m.cName }}</span>
                  <span class="cus_tag hidden-md-and-down">{{ m.area }}</span>
                  <span class="cus_remark hidden-md-and-up">{{ m.remarks }}</span>
              </a>
              <a :href="`/filmDetail?link=${m.id}`"
                 class=" content_text_tag">{{ m.name }}</a>
              <span class="cus_remark hidden-md-and-down">{{ m.remarks }}</span>
          </div>
      </div>
  </div>
</template>

<script lang="ts" setup>
defineProps({
    relateList : []
})

</script>

<style scoped>
@media (min-width: 650px) {
    .container {
        width: 100%;
        margin-top: 36px;
    }
    .title {
        padding: 0 0 5px 10px;
        text-align: left;
        border-bottom: 2px solid #777777;
    }
    .content {
        width: 100%;
        display: flex;
        flex-direction: row;
        flex-flow: wrap;
        justify-content: start;
    }
    .container .item {
        width: calc(14% - 20px);
        margin: 2px 10px 16px 10px;
    }
    .cus_content_link {
        border-radius: 5px;
        display: flex;
        padding-top: 125%;
        background-size: cover;
        overflow: hidden;
    }

    .cus_tag {
        text-align: center;
        color: rgb(255, 255, 255);
        padding: 0 3px;
        margin: 0 0 10px 8px;
        background: rgba(0, 0, 0, 0.55);
        font-size: 12px;
        border-radius: 5px;
    }

    .content_text_tag {
        display: block;
        font-size: 13px !important;
        color: rgb(221, 221, 221);
        padding: 2px 2px 2px 2px !important;
        text-align: left;
        max-width: 80%;
        max-height: 20px;
        white-space: nowrap; /* 禁止换行 */
        overflow: hidden; /* 超出部分隐藏 */
        text-overflow: ellipsis; /* 使用省略号表示被截断的部分 */
    }
    .cus_remark {
        display: block;
        color: #888888;
        width: 100%;
        font-size: 12px;
        text-align: left;
    }

}

@media (max-width: 650px) {
    .container {
        width: 100%;
        margin-top: 36px;
    }
    .title {
        padding: 0 0 5px 10px;
        text-align: left;
        border-bottom: 2px solid #777777;
    }
    .content {
        width: 100%;
        display: flex;
        flex-direction: row;
        flex-flow: wrap;
        justify-content: start;
    }
    .container .item {
        width: calc(33% - 8px);
        margin: 2px 4px 16px 4px;
    }
    .cus_content_link {
        border-radius: 5px;
        display: flex;
        padding-top: 125%;
        background-size: cover;
        overflow: hidden;
    }

    .cus_tag {
        text-align: center;
        color: rgb(255, 255, 255);
        padding: 0 3px;
        margin: 0 0 10px 8px;
        background: rgba(0, 0, 0, 0.55);
        font-size: 12px;
        border-radius: 5px;
    }

    .content_text_tag {
        font-size: 12px !important;
        color: rgb(221, 221, 221);
        width: 100%;
        padding: 2px 0 2px 0 !important;
        text-align: left;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 2;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .cus_remark {
        display: block;
        width: 100%;
        font-size: 12px;
        color: #c2c2c2;
        text-align: center;
        background: rgba(0,0,0,0.55);
        border-radius:  0 0 5px 5px;
    }

}

</style>