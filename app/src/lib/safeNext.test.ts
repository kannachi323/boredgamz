import { describe, expect, it } from "vitest";
import { safeNext } from "./safeNext";

describe("safeNext", () => {
  it("keeps same-site paths and rejects everything else", () => {
    expect(safeNext("/games/chess?join=HABCDE")).toBe("/games/chess?join=HABCDE");
    expect(safeNext("//evil.example")).toBe("/");
    expect(safeNext("https://evil.example")).toBe("/");
    expect(safeNext(null)).toBe("/");
  });
});
