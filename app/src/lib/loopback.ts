const loopbackHosts = new Set(["localhost", "127.0.0.1", "[::1]", "::1"]);

export function alignLoopbackHost(value: string, pageHostname: string) {
  try {
    const url = new URL(value);
    if (loopbackHosts.has(url.hostname) && loopbackHosts.has(pageHostname)) {
      url.hostname = pageHostname;
      return url.toString().replace(/\/$/, "");
    }
  } catch {
    // Relative paths and invalid configuration are returned unchanged.
  }
  return value;
}
