<!--
  Cover / title layout. Warm off-white paper with a soft coral wash,
  Tailscale wordmark at the top, large hero title, speaker block.
  Designed for slide 1 and workshop-level opening slides.
-->
<template>
  <div class="slidev-layout cover bg-planet-teal">
    <div class="cover-inner">
      <img
        class="cover-wordmark"
        src="../assets/tailscale-wordmark.png"
        alt="Tailscale"
      />
      <slot />
    </div>
  </div>
</template>

<style scoped>
.cover {
  padding: 3.5rem 3.5rem 3rem;
  position: relative;
}
.cover-inner {
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}
.cover-wordmark {
  height: 42px;
  width: auto;
  opacity: 0.92;
  align-self: flex-start;
}
.cover :deep(h1) {
  font-size: 3rem;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--tailscale-text-strong);
  max-width: 18ch;
  margin: 0;
}
.cover :deep(h2) {
  font-weight: 500;
  color: var(--tailscale-coral);
  font-size: 1.3rem;
  margin: 0.4rem 0 0;
}
.cover :deep(p) {
  color: var(--tailscale-text);
  font-size: 0.95rem;
  line-height: 1.7;
}
.cover :deep(strong) {
  color: var(--tailscale-text-strong);
  font-weight: 600;
}
</style>
