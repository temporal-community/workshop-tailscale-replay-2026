<!--
  Section divider. One big bold phrase centered on the Tailscale geometric bg.
  Use for "Exercise N", "Architecture", "Wrap-Up"-style breaks.

  The bg image is dark (charcoal with geometric shapes) which is the
  authentic Tailscale section-divider style, so section headings invert
  to white against that backdrop while the rest of the deck stays light.
-->
<template>
  <div class="slidev-layout section bg-planet-purple">
    <div class="section-inner">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.section {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  position: relative;
}
.section::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(31, 30, 30, 0.25) 0%, rgba(31, 30, 30, 0.55) 100%);
  pointer-events: none;
}
.section-inner {
  position: relative;
  text-align: center;
  max-width: 22ch;
}
.section :deep(h1) {
  font-size: 3.4rem;
  font-weight: 700;
  letter-spacing: -0.025em;
  color: #ffffff;
  margin: 0;
  line-height: 1.05;
}
.section :deep(h2) {
  color: var(--tailscale-coral-soft);
  font-weight: 500;
  margin-top: 0.8rem;
}
.section :deep(p) {
  color: #EDEBEA;
  margin-top: 1rem;
}
</style>
