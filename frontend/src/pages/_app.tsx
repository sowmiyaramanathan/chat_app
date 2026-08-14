import { createTheme, styled, ThemeProvider } from "@mui/material";
import type { AppProps } from "next/app";
import Navbar from "../../components/Navbar";
import { useEffect, useState } from "react";
import { useRouter } from "next/router";
import type {
  PaletteColor,
  PaletteColorOptions,
} from "@mui/material/styles";

declare module "@mui/material/styles" {
  interface Palette {
    msgBg: PaletteColor;
  }

  interface PaletteOptions {
    msgBg?: PaletteColorOptions;
  }
}

const theme = createTheme({
  palette: {
    primary: {
      //Buttons
      main: "#1f316f",
      light: "#f9dbba",
      contrastText: "#fff",
    },
    secondary: {
      //bg
      main: "#1a4870",
    },
    msgBg: {
      //msgBg
      main: "#5b99c2",
      contrastText: "#000",
    },
    // error: {
    //   main: "red",
    // },
  },
});

const BackgroundContainer = styled("div")(({ theme }) => ({
  minHeight: "100vh",
  backgroundColor: theme.palette.secondary.main,
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
