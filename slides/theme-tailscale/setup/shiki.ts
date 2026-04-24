// Register the tailscale-light Shiki theme so fenced code blocks use the deck's palette.
// (Plain default export. Slidev calls it to get the Shiki config.)
import tailscaleLight from './tailscale-light.json'

export default () => ({
  themes: {
    dark: tailscaleLight as any,
    light: tailscaleLight as any,
  },
})
