<!--
  Tailscale 3x3 dot-grid mark as inline SVG so layouts can tint it via currentColor.
  Default color is the Tailscale charcoal; pass `color` prop to override.
-->
<template>
  <svg
    viewBox="0 0 100 100"
    fill="currentColor"
    xmlns="http://www.w3.org/2000/svg"
    aria-label="Tailscale"
    role="img"
  >
    <!-- Row 1: outer dots muted, center dot filled -->
    <circle cx="15" cy="15" r="11" opacity="0.35" />
    <circle cx="50" cy="15" r="11" opacity="0.35" />
    <circle cx="85" cy="15" r="11" opacity="0.35" />
    <!-- Row 2: full intensity -->
    <circle cx="15" cy="50" r="11" />
    <circle cx="50" cy="50" r="11" />
    <circle cx="85" cy="50" r="11" />
    <!-- Row 3: outer dots muted, center filled -->
    <circle cx="15" cy="85" r="11" opacity="0.35" />
    <circle cx="50" cy="85" r="11" />
    <circle cx="85" cy="85" r="11" opacity="0.35" />
  </svg>
</template>
