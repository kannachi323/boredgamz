export type ChessTheme = 'walnut' | 'midnight' | 'sage' | 'parchment'

export const chessThemes: Record<ChessTheme, { name: string; light: string; dark: string; accent: string }> = {
  walnut: { name: 'Walnut club', light: '#d6c2a3', dark: '#6b4934', accent: '#7dcfb6' },
  midnight: { name: 'Midnight', light: '#aeb9c5', dark: '#344253', accent: '#87bff2' },
  sage: { name: 'Garden sage', light: '#d9dcc7', dark: '#72806b', accent: '#f0c77b' },
  parchment: { name: 'Parchment', light: '#eadfc8', dark: '#a88768', accent: '#8f5f4c' },
}
