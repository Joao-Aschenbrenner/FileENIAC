import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen, act } from "@testing-library/react";
import { useDocumentVisible } from "../useDocumentVisible";

// We can't simulate document.hidden reliably in jsdom default environment,
// so we test the hook's external behaviour: it should return a boolean and
// update when visibility changes are dispatched.

function Probe() {
  const visible = useDocumentVisible();
  return <div data-testid="visible">{String(visible)}</div>;
}

describe("useDocumentVisible", () => {
  beforeEach(() => {
    Object.defineProperty(document, "visibilityState", {
      value: "visible",
      writable: true,
      configurable: true,
    });
    Object.defineProperty(document, "hidden", {
      value: false,
      writable: true,
      configurable: true,
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("returns true when document is visible", () => {
    render(<Probe />);
    expect(screen.getByTestId("visible").textContent).toBe("true");
  });

  it("returns false when document is hidden", () => {
    Object.defineProperty(document, "visibilityState", { value: "hidden" });
    Object.defineProperty(document, "hidden", { value: true });

    render(<Probe />);
    expect(screen.getByTestId("visible").textContent).toBe("false");
  });

  it("adds and removes visibilitychange listener on mount/unmount", () => {
    const addSpy = vi.spyOn(document, "addEventListener");
    const removeSpy = vi.spyOn(document, "removeEventListener");

    const { unmount } = render(<Probe />);
    expect(addSpy).toHaveBeenCalledWith("visibilitychange", expect.any(Function));

    unmount();
    expect(removeSpy).toHaveBeenCalledWith("visibilitychange", expect.any(Function));
  });

  it("updates state when visibility changes from visible to hidden", () => {
    render(<Probe />);
    expect(screen.getByTestId("visible").textContent).toBe("true");

    act(() => {
      Object.defineProperty(document, "visibilityState", { value: "hidden" });
      document.dispatchEvent(new Event("visibilitychange"));
    });

    expect(screen.getByTestId("visible").textContent).toBe("false");
  });

  it("updates state when visibility changes from hidden to visible", () => {
    Object.defineProperty(document, "visibilityState", { value: "hidden" });
    Object.defineProperty(document, "hidden", { value: true });

    render(<Probe />);
    expect(screen.getByTestId("visible").textContent).toBe("false");

    act(() => {
      Object.defineProperty(document, "visibilityState", { value: "visible" });
      Object.defineProperty(document, "hidden", { value: false });
      document.dispatchEvent(new Event("visibilitychange"));
    });

    expect(screen.getByTestId("visible").textContent).toBe("true");
  });
});