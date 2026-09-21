import Head from "next/head";
import { Box, CssBaseline, GlobalStyles, ThemeProvider } from "@mui/material";
import type { AppProps } from "next/app";
import Navbar from "../../components/Navbar";
import { getTheme } from "../../components/theme";
import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import axios from "axios";
import { getTokenExpiry, setToken, setTokens } from "../../token/token";
import { STRINGS } from "../../components/keys";
import { API_BASE_URL } from "../../components/config";

const ProtectedRoutes = [
  "/user/profile",
  "/user/chats",
  "/user/discover",
  "/user/requests",
];

export default function App({ Component, pageProps }: AppProps) {
  const router = useRouter();
  const [pushed, setPushed] = useState(false);
  const [mode, setMode] = useState<"light" | "dark">("light");

  const appTheme = getTheme(mode);

  useEffect(() => {
    const savedMode = localStorage.getItem("hello-theme");
    if (savedMode === "light" || savedMode === "dark") {
      setMode(savedMode);
    }
  }, []);

  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | undefined;
    let refreshInFlight: Promise<string> | undefined;

    const expireSession = () => {
      setToken(null);
      setPushed(false);
      if (router.pathname !== "/user/signin") {
        void router.replace("/user/signin");
      }
    };

    const scheduleRefresh = (token: string) => {
      if (timer) clearTimeout(timer);
      const expiry = getTokenExpiry(token);
      if (expiry === null) {
        expireSession();
        return;
      }
      const remaining = expiry * 1000 - Date.now();
      // Use 10% of the remaining token lifetime as the refresh buffer (up to
      // one minute). A fixed one-minute buffer schedules immediately when
      // testing with one-minute access tokens, causing a refresh loop.
      const refreshBuffer = Math.min(60_000, Math.max(1_000, remaining * 0.1));
      const delay = Math.max(1_000, remaining - refreshBuffer);
      timer = setTimeout(() => {
        void refreshAccessToken().catch(expireSession);
      }, delay);
    };

    const refreshAccessToken = async (): Promise<string> => {
      if (refreshInFlight) return refreshInFlight;

      const refreshToken = localStorage.getItem("refreshToken");
      if (!refreshToken) throw new Error("No refresh token");

      refreshInFlight = axios
        .post(`${API_BASE_URL}/user/auth/refresh`, { refreshToken })
        .then((response) => {
          const { token, refreshToken: nextRefreshToken } = response.data;
          if (
            typeof token !== "string" ||
            typeof nextRefreshToken !== "string"
          ) {
            throw new Error("Invalid refresh response");
          }
          setTokens(token, nextRefreshToken);
          scheduleRefresh(token);
          return token;
        })
        .finally(() => {
          refreshInFlight = undefined;
        });
      return refreshInFlight;
    };

    const interceptor = axios.interceptors.response.use(
      (response) => response,
      async (error) => {
        const request = error.config as
          | (typeof error.config & { _retry?: boolean })
          | undefined;
        const isRefreshRequest =
          request?.url === `${API_BASE_URL}/user/auth/refresh`;
        if (
          error.response?.status !== 401 ||
          !request ||
          request._retry ||
          isRefreshRequest
        ) {
          if (isRefreshRequest) expireSession();
          return Promise.reject(error);
        }

        try {
          request._retry = true;
          const token = await refreshAccessToken();
          request.headers = request.headers ?? {};
          request.headers.Authorization = `Bearer ${token}`;
          return axios(request);
        } catch {
          expireSession();
        }
        return Promise.reject(error);
      },
    );

    const token = localStorage.getItem("token");
    if (token) scheduleRefresh(token);

    return () => {
      axios.interceptors.response.eject(interceptor);
      if (timer) clearTimeout(timer);
    };
  }, [router]);

  const toggleMode = () => {
    setMode((currentMode) => {
      const nextMode = currentMode === "light" ? "dark" : "light";
      localStorage.setItem("hello-theme", nextMode);
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
    <>
      <Head>
        <title>{STRINGS.app.title}</title>
        <meta name="keywords" content={STRINGS.app.metaKeywords} />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

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
    </>
  );
}
