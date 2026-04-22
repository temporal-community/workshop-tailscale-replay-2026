<!--
  Workshop table of contents. Rendered once before each section divider
  with `current` set to the id of the section we're about to enter.
  Active row is highlighted in mint; the rest read as muted Inter Light.
-->
<script setup lang="ts">
defineProps<{ current?: string }>()

const sections = [
  { id: 'arch',       label: 'Architecture' },
  { id: 'ex1',        label: 'Exercise 1: Hello Tailnet' },
  { id: 'ex2',        label: 'Exercise 2: Explore Tailscale' },
  { id: 'agents',     label: 'AI Agents on Temporal' },
  { id: 'ex3',        label: 'Exercise 3: Weather Agent' },
  { id: 'ratelimit',  label: 'Rate Limit Demo' },
  { id: 'tsnet',      label: 'temporal-ts-net and Metrics Watcher' },
  { id: 'wrap',       label: 'Wrap-Up' },
]
</script>

<template>
  <ol class="workshop-toc">
    <li
      v-for="(section, index) in sections"
      :key="section.id"
      :class="{ active: current === section.id }"
    >
      <span class="number">{{ String(index + 1).padStart(2, '0') }}</span>
      <span class="label">{{ section.label }}</span>
    </li>
  </ol>
</template>

<style scoped>
.workshop-toc {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
}
.workshop-toc li {
  display: flex;
  align-items: baseline;
  gap: 1.25rem;
  padding: 0.1rem 0;
  color: var(--temporal-text-muted);
  font-weight: 300;
  font-size: 1.5rem;
  transition: color 150ms ease;
}
.workshop-toc li .number {
  font-variant-numeric: tabular-nums;
  font-size: 1.05rem;
  letter-spacing: 0.1em;
  color: var(--temporal-text-muted);
  opacity: 0.65;
}
.workshop-toc li.active {
  color: #ffffff;
}
.workshop-toc li.active .number {
  color: var(--temporal-green);
  opacity: 1;
}
.workshop-toc li.active .label {
  font-weight: 400;
}
</style>
