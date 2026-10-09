<script setup lang="ts">
// 休息图标：主显示器右下角的一个小笑脸（窗口见 window_options.go breakIconOptions）。
// 窗口随休息创建、休息结束时销毁，所以这里不需要任何后端状态，也不接收 tick 事件。
// 进入后持续眨眼（周期性眨两下再停顿），直到休息结束窗口被撤掉。
import { onMounted } from 'vue'

onMounted(() => {
  // 窗口透明（同 PreBreakNotice），把 html/body 的底色拉回透明，只显示图标本身。
  document.documentElement.style.background = 'transparent'
  document.body.style.background = 'transparent'
})
</script>

<template>
  <div class="break-icon">
    <svg viewBox="0 0 64 64" class="face" aria-hidden="true">
      <circle cx="32" cy="32" r="28" class="face__bg" />
      <ellipse cx="22" cy="26" rx="3.4" ry="4.6" class="eye" />
      <ellipse cx="42" cy="26" rx="3.4" ry="4.6" class="eye" />
      <path d="M19 38 Q32 50 45 38" class="mouth" />
    </svg>
  </div>
</template>

<style scoped>
.break-icon {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  user-select: none;
  -webkit-user-select: none;
  cursor: default;
}

.face {
  width: 84px;
  height: 84px;
  filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.25));
  animation: pm-icon-in 0.4s cubic-bezier(0.22, 1, 0.36, 1) both;
}
.face__bg {
  fill: #f6c64f;
}
.eye {
  fill: #3a3126;
  transform-box: fill-box;
  transform-origin: center;
  /* 0.3s 后开始，无限循环：每个周期连续眨两下，之后停顿到周期结束再眨 */
  animation: pm-blink 3s ease-in-out 0.3s infinite;
}
.mouth {
  fill: none;
  stroke: #3a3126;
  stroke-width: 3.5;
  stroke-linecap: round;
}

@keyframes pm-icon-in {
  from { opacity: 0; transform: scale(0.6); }
  to   { opacity: 1; transform: scale(1); }
}
@keyframes pm-blink {
  0%   { transform: scaleY(1); }
  4%   { transform: scaleY(0.1); }
  8%   { transform: scaleY(1); }
  12%  { transform: scaleY(0.1); }
  16%, 100% { transform: scaleY(1); }
}
</style>
