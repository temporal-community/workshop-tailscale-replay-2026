<!--
  Cover / title layout. Teal-planet background, large hero title, speaker block.
  Designed for slide 1 and workshop-level opening slides.
-->
<script setup lang="ts">
import TemporalLogo from '../components/TemporalLogo.vue'
</script>

<template>
  <div class="slidev-layout cover bg-planet-teal">
    <div class="cover-inner">
      <TemporalLogo class="cover-logo" />
      <slot />
    </div>
  </div>
</template>

<style scoped>
.cover {
  padding: 3.5rem 3.5rem 3rem;
  position: relative;
}
.cover::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(11, 16, 32, 0.15) 0%, rgba(11, 16, 32, 0.55) 100%);
  pointer-events: none;
}
.cover-inner {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.cover-logo {
  width: 44px;
  height: 44px;
  color: var(--temporal-green);
}
.cover :deep(h1) {
  font-size: 3rem;
  font-weight: 200;
  letter-spacing: -0.02em;
  color: #ffffff;
  max-width: 18ch;
  margin: 0;
}
.cover :deep(h2) {
  font-weight: 300;
  color: var(--temporal-lavender);
  font-size: 1.3rem;
  margin: 0.4rem 0 0;
}
.cover :deep(p) {
  color: var(--temporal-text);
  font-size: 0.95rem;
  line-height: 1.7;
}
.cover :deep(strong) {
  color: #ffffff;
  font-weight: 500;
}
</style>
