export function seriesLabel(bestOf: number) {
  return bestOf > 1 ? `Best of ${bestOf}` : "";
}

export function clockLabel(control: string) {
  return !control || control === "untimed" ? "" : `${control} per move`;
}
