import { useEffect, useRef, useState } from "react";

export function useLocalTimer(serverSeconds: number) {
  const [timer, setTimer] = useState(serverSeconds);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Sync with server's time
    setTimer(serverSeconds);

    // Clear existing interval if re-syncing
    if (intervalRef.current) clearInterval(intervalRef.current);

    intervalRef.current = setInterval(() => {
      setTimer((prev) => Math.max(prev - 1, 0));
    }, 1000);

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current);
    };
  }, [serverSeconds]);

  return timer;
}
