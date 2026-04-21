// Default Mermaid config for the Temporal Slidev theme.
// Dark base with Temporal palette overrides — mint borders, lavender edges,
// white text on deep-slate nodes.
// (Plain default export — Slidev calls it to get the config.)

export default () => ({
  theme: 'base',
  themeVariables: {
    fontFamily: 'Inter, system-ui, sans-serif',
    fontSize: '18px',
    background: '#0B1020',
    primaryColor: '#111827',
    primaryTextColor: '#F8FAFC',
    primaryBorderColor: '#59FDA0',
    secondaryColor: '#1a2338',
    secondaryTextColor: '#C7D7FE',
    secondaryBorderColor: '#7C3AED',
    tertiaryColor: '#0B1020',
    tertiaryTextColor: '#C3A0FF',
    tertiaryBorderColor: '#7C8FB1',
    lineColor: '#C3A0FF',
    textColor: '#C7D7FE',
    mainBkg: '#111827',
    nodeBorder: '#59FDA0',
    clusterBkg: 'rgba(124, 58, 237, 0.10)',
    clusterBorder: '#7C3AED',
    edgeLabelBackground: '#0B1020',
    labelBackground: '#0B1020',
    actorBkg: '#111827',
    actorBorder: '#59FDA0',
    actorTextColor: '#F8FAFC',
    actorLineColor: '#C3A0FF',
    signalColor: '#C7D7FE',
    signalTextColor: '#C7D7FE',
    noteBkgColor: 'rgba(124, 58, 237, 0.14)',
    noteTextColor: '#C7D7FE',
    noteBorderColor: '#7C3AED',
  },
})
