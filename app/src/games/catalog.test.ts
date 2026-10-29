import { describe, expect, it } from "vitest";
import { gameForCode, games, makeInviteCode, normalizeCode } from "./catalog";

describe("private codes", () => {
  it("encode their game so a code alone finds the table", () => {
    for (const game of games) {
      const code = makeInviteCode(game);
      expect(code).toMatch(/^[A-Z0-9]{6}$/);
      expect(gameForCode(code)?.slug).toBe(game.slug);
    }
  });

  it("normalises user input and rejects unknown prefixes", () => {
    expect(normalizeCode(" h-ab c12 ")).toBe("HABC12");
    expect(gameForCode("habc12")?.slug).toBe("chess");
    expect(gameForCode("ZABC12")).toBeUndefined();
    expect(gameForCode("HAB")).toBeUndefined();
  });

  it("uses unique prefixes", () => {
    expect(new Set(games.map((game) => game.codePrefix)).size).toBe(games.length);
  });
});
