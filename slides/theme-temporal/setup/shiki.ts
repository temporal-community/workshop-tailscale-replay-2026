// Register the temporal-dark Shiki theme so fenced code blocks use the deck's palette.
// (Plain default export. Slidev calls it to get the Shiki config.)
import temporalDark from './temporal-dark.json'

export default () => ({
  themes: {
    dark: temporalDark as any,
    light: temporalDark as any,
  },
})
