<!--
  Two-column layout with a full-width header row above both columns.
  Markdown:
    ---
    layout: two-cols-header
    ---
    # Title spanning both columns
    ::left::
    Left column content
    ::right::
    Right column content
-->
<script setup lang="ts">
import TailscaleFooter from '../components/TailscaleFooter.vue'
</script>

<template>
  <div class="slidev-layout two-cols-header bg-grid">
    <div class="header"><slot name="default" /></div>
    <div class="two-cols-grid">
      <div class="col col-left"><slot name="left" /></div>
      <div class="col col-right"><slot name="right" /></div>
    </div>
    <TailscaleFooter />
  </div>
</template>

<style scoped>
.two-cols-header {
  display: flex;
  flex-direction: column;
  height: 100%;
}
.header {
  flex: 0 0 auto;
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
.two-cols-header :deep(.col) p,
.two-cols-header :deep(.col) li {
  font-size: 1.3rem;
  line-height: 1.35;
}
</style>
