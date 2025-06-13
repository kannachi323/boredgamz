import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

import { AuthContext } from "./AuthContext";



export const AuthProvider = ({ children } : { children: React.ReactNode }) => {
  const navigate = useNavigate();

  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [user, setUser] = useState<{ id: string; username: string } | null>(null);

  useEffect(() => {
    async function checkAuth() {
      try {
        const res = await fetch("http://localhost:3000/auth", {
          method: "GET",
          credentials: "include",
        });

        if (res.ok) {
          const data = await res.json();
          setUser({ id: data.id, username: data.username });
          setIsAuthenticated(true);
        } else {
          setIsAuthenticated(false);
          navigate("/signup");
        }
      } catch (err) {
        console.error("Auth check failed:", err);
        setIsAuthenticated(false);
        navigate("/signup");
      }
    }

    checkAuth();
  }, [navigate]);


  return (
    <AuthContext.Provider value={{ isAuthenticated, setIsAuthenticated, user, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}