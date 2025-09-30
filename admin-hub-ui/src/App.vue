<script setup lang="ts">
import { useThemeStore } from '@/stores/modules/theme.ts'
import variables from '@/styles/variable.module.scss'

const app = useThemeStore()

// 切换到下一个主题色
function switchToNextTheme() {
  const themeKeys = Object.keys(variables)
  const currentIndex = themeKeys.indexOf(app.themeName)
  const nextIndex = (currentIndex + 1) % themeKeys.length
  app.setThemeName(themeKeys[nextIndex])
}
</script>

<template>
  <a-config-provider :theme="app.themeConfig">
    <a-select v-model:value="app.themeName" style="width: 240px">
      <a-select-option v-for="(color, name) in variables" :value="name">
        <span
          :style="{
            display: 'inline-block',
            width: '16px',
            height: '16px',
            backgroundColor: color,
            marginRight: '8px',
            borderRadius: '4px',
          }"
        ></span>
        {{ name }}:{{ color }}
      </a-select-option>
    </a-select>
    <a-select v-model:value="app.darkMode" style="width: 120px">
      <a-select-option value="dark">dark</a-select-option>
      <a-select-option value="light">light</a-select-option>
    </a-select>
    <a-button-group>
      <a-button @click="switchToNextTheme" type="primary"
          >切换主题-
          {{ app.themeName }}
        </a-button>
      <a-button @click="app.toggleDarkMode">切换模式{{ app.darkModeComp }}</a-button>
    </a-button-group>
    <div class="test">test</div>
  </a-config-provider>
</template>

<style lang="scss">
@use '@/styles/theme.scss';

.test {
  // 直接从store中获取主题色，这样刷新后也能保持
  background: v-bind('variables[app.themeName]');
}
</style>
