<!--
  Two-column layout optimised for a diagram/figure in the right column.
  Left column: 40% (text/bullets). Right column: 60% (image, vertically centered).
  Markdown:
    ---
    layout: two-cols-figure
    ---
    ## Left column
    ::right::
    <img … class="tailscale-figure" />
-->
<script setup lang="ts">
import TailscaleFooter from '../components/TailscaleFooter.vue'
</script>

<template>
  <div class="slidev-layout two-cols-figure bg-grid">
    <div class="two-cols-grid">
      <div class="col col-left"><slot /></div>
      <div class="col col-right"><slot name="right" /></div>
    </div>
    <TailscaleFooter />
  </div>
</template>

<style scoped>
.two-cols-figure {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.two-cols-grid {
  flex: 1 1 auto;
  min-height: 0;
  display: grid;
  grid-template-columns: 2fr 3fr;
  gap: 2rem;
  align-items: stretch;
  overflow: hidden;
}
.two-cols-figure :deep(h1) {
  font-size: 2.2rem;
  margin-bottom: 0.6rem;
}
.two-cols-figure :deep(.col) p,
.two-cols-figure :deep(.col) li {
  font-size: 1.3rem;
  line-height: 1.35;
}
.two-cols-figure :deep(.col-right) {
  display: flex;
  align-items: center;
  justify-content: center;
}
.two-cols-figure :deep(.col-right) img {
  width: 100%;
  height: auto;
  max-height: 100%;
  max-width: 100%;
}
</style>
