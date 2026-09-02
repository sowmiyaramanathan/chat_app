import { Box, CssBaseline, GlobalStyles, ThemeProvider } from "@mui/material";
import type { AppProps } from "next/app";
import Navbar from "../../components/Navbar";
import { getTheme } from "../../components/theme";
import { useEffect, useState } from "react";
import { useRouter } from "next/router";

const ProtectedRoutes = ["/user/profile", "/user/chats", "/user/requests"];

export default function App({ Component, pageProps }: AppProps) {
  const router = useRouter();
  const [pushed, setPushed] = useState(false);
  const [mode, setMode] = useState<"light" | "dark">("light");

  const appTheme = getTheme(mode);

  useEffect(() => {
    const savedMode = localStorage.getItem("whisper-theme");
    if (savedMode === "light" || savedMode === "dark") {
      setMode(savedMode);
    }
  }, []);

  const toggleMode = () => {
    setMode((currentMode) => {
      const nextMode = currentMode === "light" ? "dark" : "light";
      localStorage.setItem("whisper-theme", nextMode);
      return nextMode;
    });
  };

  useEffect(() => {
    const token = localStorage.getItem("token");
    const isProtectedRoute = ProtectedRoutes.includes(router.pathname);

    if (pushed) {
      return;
    }
    if (isProtectedRoute && !token) {
      router.push("/user/signin");
      setPushed(true);
    } else if (
      (router.pathname === "/user/signin" ||
        router.pathname === "/user/signup") &&
      token
    ) {
      router.push("/user/profile");
      setPushed(true);
    }
  }, [router, pushed]);

  return (
    <ThemeProvider theme={appTheme}>
      <CssBaseline />
      <GlobalStyles
        styles={{
          "@keyframes floatIn": {
            from: { opacity: 0, transform: "translateY(10px)" },
            to: { opacity: 1, transform: "translateY(0)" },
          },
        }}
      />
      <Box
        sx={{
          minHeight: "100vh",
          position: "relative",
          overflow: "hidden",
          background: (theme) =>
            `linear-gradient(140deg, ${theme.palette.background.default} 0%, ${theme.palette.msgBg.light} 100%)`,
          "&::before": {
            content: '""',
            position: "fixed",
            inset: 0,
            opacity: 0.18,
            pointerEvents: "none",
            backgroundImage: (theme) =>
              `repeating-linear-gradient(17deg, transparent 0 22px, ${theme.palette.primary.main} 23px 24px, transparent 25px 44px)`,
            maskImage: "linear-gradient(to bottom, black, transparent 75%)",
          },
          "&::after": {
            content: '""',
            position: "fixed",
            width: 360,
            height: 360,
            right: -130,
            bottom: -150,
            borderRadius: "50%",
            background: (theme) => theme.palette.secondary.light,
            opacity: 0.38,
            pointerEvents: "none",
          },
        }}
      >
        <Box sx={{ position: "relative", zIndex: 1 }}>
          <Navbar mode={mode} onToggleMode={toggleMode} />
          <Component {...pageProps} />
        </Box>
      </Box>
    </ThemeProvider>
  );
}
