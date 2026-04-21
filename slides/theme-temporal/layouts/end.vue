<!--
  Closing layout — Rex constellation background, centered sign-off.
  Use for the final "Questions?" / "Thank you" slide.
-->
<script setup lang="ts">
import TemporalLogo from '../components/TemporalLogo.vue'
</script>

<template>
  <div class="slidev-layout end bg-rex">
    <div class="end-inner">
      <TemporalLogo class="end-logo" />
      <div class="end-content">
        <slot />
      </div>
    </div>
  </div>
</template>

<style scoped>
.end {
  padding: 3.5rem;
  position: relative;
}
.end::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(11, 16, 32, 0.2) 0%, rgba(11, 16, 32, 0.6) 100%);
  pointer-events: none;
}
.end-inner {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.end-logo {
  width: 38px;
  height: 38px;
  color: var(--temporal-green);
}
.end-content {
  align-self: flex-start;
  max-width: 28ch;
}
.end :deep(h1) {
  font-size: 3.2rem;
  font-weight: 200;
  letter-spacing: -0.02em;
  color: #ffffff;
  margin: 0 0 1.4rem;
}
.end :deep(p) {
  color: var(--temporal-text);
  line-height: 1.7;
}
.end :deep(strong) {
  color: #ffffff;
}
.end :deep(a) {
  color: var(--temporal-green);
}
</style>
