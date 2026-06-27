import { useEffect, useRef, useState } from "react";

export function useDocumentVisible(): boolean {
  const [visible, setVisible] = useState(!document.hidden);
  const visibleRef = useRef(true);

  useEffect(() => {
    visibleRef.current = !document.hidden;
    const onChange = () => {
      visibleRef.current = document.visibilityState === "visible";
      setVisible(visibleRef.current);
    };
    document.addEventListener("visibilitychange", onChange);
    return () => document.removeEventListener("visibilitychange", onChange);
  }, []);

  return visible;
}