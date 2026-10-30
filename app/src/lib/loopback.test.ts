import { describe, expect, it } from "vitest";
import { alignLoopbackHost } from "./loopback";

describe("loopback endpoint alignment", () => {
  it("keeps development cookies and sockets on the page hostname", () => {
    expect(alignLoopbackHost("http://localhost:3000/api", "127.0.0.1")).toBe("http://127.0.0.1:3000/api");
    expect(alignLoopbackHost("ws://localhost:3000/api/v1/arena", "127.0.0.1")).toBe("ws://127.0.0.1:3000/api/v1/arena");
  });

  it("does not rewrite production hosts", () => {
    expect(alignLoopbackHost("https://api.example.com", "localhost")).toBe("https://api.example.com");
  });
});
