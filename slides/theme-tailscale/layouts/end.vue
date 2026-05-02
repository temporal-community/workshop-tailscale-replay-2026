<!--
  Closing layout. Soft coral glow background, centered sign-off with Tailscale mark.
  Use for the final "Questions?" / "Thank you" slide.
-->
<script setup lang="ts">
import TailscaleLogo from '../components/TailscaleLogo.vue'
</script>

<template>
  <div class="slidev-layout end bg-glow">
    <div class="end-inner">
      <TailscaleLogo class="end-logo" />
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
.end-inner {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.end-logo {
  width: 44px;
  height: 44px;
  color: var(--tailscale-coral);
}
.end-content {
  align-self: flex-start;
  max-width: 22ch;
}
.end :deep(h1) {
  font-size: 3.2rem;
  font-weight: 700;
  letter-spacing: -0.025em;
  color: var(--tailscale-text-strong);
  margin: 0 0 1.4rem;
}
.end :deep(p) {
  color: var(--tailscale-text);
  line-height: 1.7;
}
.end :deep(strong) {
  color: var(--tailscale-text-strong);
}
.end :deep(a) {
  color: var(--tailscale-coral);
}
</style>
