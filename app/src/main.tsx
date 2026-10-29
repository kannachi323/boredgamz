import "@fontsource-variable/geist";
import React from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import "./index.css";
import { applyTheme, readTheme } from "./lib/theme";
import router from "./router";

applyTheme(readTheme());

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
