<!--
  Footer row shown on content layouts — workshop title left, slide index right.
  Positioned via flex (margin-top: auto) rather than absolute so it never overlaps.
-->
<script setup lang="ts">
import { useNav } from '@slidev/client'
import { computed } from 'vue'

const { currentPage, total } = useNav()
const pageLabel = computed(() =>
  `${String(currentPage.value).padStart(2, '0')} / ${String(total.value).padStart(2, '0')}`
)
</script>

<template>
  <div class="temporal-footer">
    <span>Replay 2026 · Tailscale × Temporal</span>
    <span class="page">{{ pageLabel }}</span>
  </div>
</template>

<style scoped>
.page {
  font-variant-numeric: tabular-nums;
  color: var(--temporal-green);
  opacity: 0.9;
}
</style>
