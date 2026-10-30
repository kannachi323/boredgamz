import { describe, expect, it } from "vitest";
import { gameForCode, games, makeInviteCode, normalizeCode } from "./catalog";

describe("private codes", () => {
  it("are six digits that encode their game", () => {
    for (const game of games) {
      const code = makeInviteCode(game);
      expect(code).toMatch(/^\d{6}$/);
      expect(gameForCode(code)?.slug).toBe(game.slug);
    }
  });

  it("normalise user input and reject unknown prefixes", () => {
    expect(normalizeCode(" 412-345 ")).toBe("412345");
    expect(gameForCode("412345")?.slug).toBe("chess");
    expect(gameForCode("912345")).toBeUndefined();
    expect(gameForCode("412")).toBeUndefined();
  });

  it("use unique prefixes", () => {
    expect(new Set(games.map((game) => game.codePrefix)).size).toBe(games.length);
  });
});
