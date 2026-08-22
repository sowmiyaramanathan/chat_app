import { styled, ThemeProvider } from "@mui/material";
import type { AppProps } from "next/app";
import Navbar from "../../components/Navbar";
import { theme } from "../../components/theme";
import { useEffect, useState } from "react";
import { useRouter } from "next/router";

const BackgroundContainer = styled("div")(({ theme }) => ({
  minHeight: "100vh",
  background: `linear-gradient(160deg, ${theme.palette.secondary.light} 0%, ${theme.palette.background.default} 45%, ${theme.palette.msgBg.light} 100%)`,
}));

const ProtectedRoutes = ["/user/profile", "/user/chats", "/user/requests"];

export default function App({ Component, pageProps }: AppProps) {
  const router = useRouter();
  const [pushed, setPushed] = useState(false);

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
    <ThemeProvider theme={theme}>
      <BackgroundContainer>
        <Navbar />
        <Component {...pageProps} />
      </BackgroundContainer>
    </ThemeProvider>
  );
}
