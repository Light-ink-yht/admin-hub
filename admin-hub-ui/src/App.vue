<script setup lang="ts">
import {useAppStore} from "@/stores/modules/theme.ts"
import variables from '@/styles/variable.module.scss'
import {ref} from "vue";

const app = useAppStore()
const backgroundColor = ref('red')

function toggleColor(color: string) {
  console.log(color)
  backgroundColor.value = color
}
</script>

<template>
  <a-config-provider :theme="app.themeConfig">
    <a-select v-model:value="app.themeName" style="width: 240px">
      <a-select-option v-for="(color, name) in variables" :value="name"> {{ name }}:{{ color }}</a-select-option>
    </a-select>
    <a-select v-model:value="app.darkMode" style="width: 120px">
      <a-select-option value="dark">dark</a-select-option>
      <a-select-option value="light">light</a-select-option>
    </a-select>
    <a-button-group>
      <a-button  @click="toggleColor(variables[app.themeName])" type="primary">切换主题- {{ app.themeName }}</a-button>
      <a-button @click="app.toggleDarkMode">切换模式{{ app.darkModeComp }}</a-button>
    </a-button-group>
    <div class="test">test</div>
  </a-config-provider>
</template>

<style lang="scss">
@use "@/styles/theme.scss";

.test {
  background: v-bind(backgroundColor);
}
</style>
