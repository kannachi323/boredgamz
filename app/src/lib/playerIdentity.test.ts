import { describe, expect, it } from "vitest";
import { combinePlayerIdentity } from "./playerIdentity";

describe("player connection identity", () => {
  it("gives separate browser tabs distinct multiplayer identities", () => {
    const browserID = "6860a6ed-ed94-4983-b229-978cbdcb01a1";

    const firstTab = combinePlayerIdentity(browserID, "tab-one");
    const secondTab = combinePlayerIdentity(browserID, "tab-two");

    expect(firstTab).not.toBe(secondTab);
    expect(firstTab).toBe(`${browserID}:tab-one`);
    expect(secondTab).toBe(`${browserID}:tab-two`);
  });

  it("keeps the identifier within the server's bounded player ID limit", () => {
    const playerID = combinePlayerIdentity("x".repeat(100), "y".repeat(100));

    expect(`${playerID}:spectator`.length).toBeLessThanOrEqual(80);
  });
});
