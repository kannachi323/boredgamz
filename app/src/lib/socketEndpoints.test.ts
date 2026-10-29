import { describe, expect, it } from "vitest";
import { resolveSocketRoot } from "./socketEndpoints";

describe("socket endpoint configuration", () => {
  it("uses the documented Arena root for Chess and sibling services", () => {
    const config = {
      arenaRoot: "ws://games.test/api/v1/arena/",
      protocol: "http:",
      host: "localhost:5173",
    };

    expect(resolveSocketRoot("arena", config)).toBe(
      "ws://games.test/api/v1/arena",
    );
    expect(resolveSocketRoot("poker", config)).toBe(
      "ws://games.test/api/v1/poker",
    );
    expect(resolveSocketRoot("mahjong", config)).toBe(
      "ws://games.test/api/v1/mahjong",
    );
  });

  it("supports the repository's legacy WebSocket API root", () => {
    const config = {
      legacyApiRoot: "ws://localhost:3000/api",
      protocol: "http:",
      host: "localhost:5173",
    };

    expect(resolveSocketRoot("arena", config)).toBe(
      "ws://localhost:3000/api/v1/arena",
    );
    expect(resolveSocketRoot("poker", config)).toBe(
      "ws://localhost:3000/api/v1/poker",
    );
  });

  it("falls back to same-origin secure sockets in production", () => {
    expect(
      resolveSocketRoot("arena", {
        protocol: "https:",
        host: "play.example.com",
      }),
    ).toBe("wss://play.example.com/api/v1/arena");
  });
});
