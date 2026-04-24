<!--
  Two-column layout with Tailscale spacing.
  Markdown:
    ---
    layout: two-cols
    ---
    ## Left column
    ::right::
    ## Right column
-->
<script setup lang="ts">
import TailscaleFooter from '../components/TailscaleFooter.vue'
</script>

<template>
  <div class="slidev-layout two-cols bg-grid">
    <div class="two-cols-grid">
      <div class="col col-left"><slot /></div>
      <div class="col col-right"><slot name="right" /></div>
    </div>
    <TailscaleFooter />
  </div>
</template>

<style scoped>
.two-cols {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.two-cols-grid {
  flex: 1 1 auto;
  min-height: 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2.5rem;
  align-items: start;
  overflow: hidden;
}
/* Two-cols slides are cramped when the h1 wraps to 2-3 lines. Shrink
   the h1 here so it fits on 1-2 lines and leaves room for body content. */
.two-cols :deep(h1) {
  font-size: 2.2rem;
  margin-bottom: 0.6rem;
}
.two-cols :deep(.col) p,
.two-cols :deep(.col) li {
  font-size: 1.3rem;
  line-height: 1.35;
}
</style>
