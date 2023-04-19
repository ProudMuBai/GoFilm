<template>
    <div class="relate_container">
        <h2 class="title">相关推荐</h2>
        <FilmList :list="relateList"/>
    </div>
</template>

<script lang="ts" setup>
import FilmList from "./FilmList.vue";

defineProps({
    relateList: Array
})

</script>

<style scoped>
.relate_container {
    width: 100%;
    margin-top: 36px;
}

.title {
    padding: 0 0 5px 10px;
    text-align: left;
    border-bottom: 2px solid #777777;
}

</style>