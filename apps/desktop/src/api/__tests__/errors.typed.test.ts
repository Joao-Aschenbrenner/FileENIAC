import { describe, it, expect } from "vitest";
import { ApiError } from "../errors";

describe("ApiError class", () => {
  it("new ApiError(401, url, msg) sets status, url, message, name", () => {
    const err = new ApiError(401, "http://localhost:8080/api/sessions", "token expired");
    expect(err.status).toBe(401);
    expect(err.url).toBe("http://localhost:8080/api/sessions");
    expect(err.message).toBe("token expired");
    expect(err.name).toBe("ApiError");
  });

  it("isUnauthorized() returns true only for 401", () => {
    expect(new ApiError(401, "u", "m").isUnauthorized()).toBe(true);
    expect(new ApiError(403, "u", "m").isUnauthorized()).toBe(false);
    expect(new ApiError(500, "u", "m").isUnauthorized()).toBe(false);
    expect(new ApiError(200, "u", "m").isUnauthorized()).toBe(false);
  });

  it("isForbidden() returns true only for 403", () => {
    expect(new ApiError(403, "u", "m").isForbidden()).toBe(true);
    expect(new ApiError(401, "u", "m").isForbidden()).toBe(false);
    expect(new ApiError(500, "u", "m").isForbidden()).toBe(false);
    expect(new ApiError(404, "u", "m").isForbidden()).toBe(false);
  });

  it("isTimeout() always false (TimeoutError is a separate class)", () => {
    expect(new ApiError(401, "u", "m").isTimeout()).toBe(false);
    expect(new ApiError(500, "u", "m").isTimeout()).toBe(false);
    expect(new ApiError(408, "u", "m").isTimeout()).toBe(false);
  });

  it("instanceof Error works (so try/catch in Effect handlers catches it)", () => {
    const err = new ApiError(500, "u", "boom");
    expect(err).toBeInstanceOf(Error);
    expect(err).toBeInstanceOf(ApiError);

    let caught: unknown;
    try {
      throw err;
    } catch (e) {
      caught = e;
    }
    expect(caught).toBeInstanceOf(Error);
    expect(caught).toBeInstanceOf(ApiError);
    expect((caught as ApiError).status).toBe(500);
  });
});
