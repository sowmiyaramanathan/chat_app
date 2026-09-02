import { createTheme } from "@mui/material/styles";
import type { PaletteColor, PaletteColorOptions } from "@mui/material/styles";

declare module "@mui/material/styles" {
  interface Palette {
    msgBg: PaletteColor;
    surface: PaletteColor;
    accent: PaletteColor;
  }

  interface PaletteOptions {
    msgBg?: PaletteColorOptions;
    surface?: PaletteColorOptions;
    accent?: PaletteColorOptions;
  }
}

const lightPalette = {
  primary: {
    main: "#18324B",
    light: "#315B72",
    dark: "#102235",
    contrastText: "#FFF9EC",
  },
  secondary: {
    main: "#F6C453",
    light: "#FFE7A3",
    dark: "#C89320",
    contrastText: "#18324B",
  },
  msgBg: {
    main: "#DDF1E6",
    light: "#F0FAF3",
    dark: "#B5D9C2",
    contrastText: "#18324B",
  },
  surface: {
    main: "#FFF9EC",
    light: "#FFFEF8",
    dark: "#F4EBD8",
    contrastText: "#18324B",
  },
  accent: {
    main: "#F07863",
    light: "#FFD0C4",
    dark: "#C95045",
    contrastText: "#FFF9EC",
  },
  background: { default: "#F7EFD9", paper: "#FFF9EC" },
  text: { primary: "#18324B", secondary: "#61707A" },
  error: { main: "#C95045" },
  success: { main: "#3E8661" },
};

const darkPalette = {
  primary: {
    main: "#F6C453",
    light: "#FFE7A3",
    dark: "#C89320",
    contrastText: "#182333",
  },
  secondary: {
    main: "#315B72",
    light: "#26465B",
    dark: "#193143",
    contrastText: "#FFF9EC",
  },
  msgBg: {
    main: "#294D48",
    light: "#356259",
    dark: "#1F3B38",
    contrastText: "#F4F0E5",
  },
  surface: {
    main: "#1B2A3A",
    light: "#23374A",
    dark: "#14212E",
    contrastText: "#F4F0E5",
  },
  accent: {
    main: "#FF8B70",
    light: "#6B3D3B",
    dark: "#D96655",
    contrastText: "#182333",
  },
  background: { default: "#101A27", paper: "#1B2A3A" },
  text: { primary: "#F4F0E5", secondary: "#B7C2C5" },
  error: { main: "#FF8B70" },
  success: { main: "#76C69A" },
};

export function getTheme(mode: "light" | "dark") {
  const palette = mode === "dark" ? darkPalette : lightPalette;

  return createTheme({
    palette: { mode, ...palette },
    typography: {
      fontFamily: '"Comic Neue", "Trebuchet MS", sans-serif',
      h1: { fontWeight: 700, letterSpacing: "-0.035em" },
      h2: { fontWeight: 700, letterSpacing: "-0.03em" },
      h3: { fontWeight: 700, letterSpacing: "-0.025em" },
      h4: { fontWeight: 700, letterSpacing: "-0.02em" },
      h5: { fontWeight: 700 },
      h6: { fontWeight: 700 },
      button: { fontWeight: 700, textTransform: "none" as const },
    },
    shape: { borderRadius: 22 },
    components: {
      MuiCssBaseline: {
        styleOverrides: {
          "*": { boxSizing: "border-box" },
          body: {
            margin: 0,
            transition: "background-color 180ms ease, color 180ms ease",
          },
          "::selection": {
            backgroundColor: palette.accent.main,
            color: palette.accent.contrastText,
          },
        },
      },
      MuiButton: {
        styleOverrides: {
          root: {
            borderRadius: 14,
            padding: "9px 19px",
            boxShadow: "none",
            transition:
              "transform 160ms ease, box-shadow 160ms ease, background-color 160ms ease",
            "&:hover": {
              transform: "translateY(-2px)",
              boxShadow: "0 7px 0 rgba(24, 50, 75, 0.14)",
            },
          },
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: {
            borderRadius: 14,
            transition: "transform 160ms ease, background-color 160ms ease",
            "&:hover": { transform: "rotate(-5deg) scale(1.05)" },
          },
        },
      },
      MuiPaper: { styleOverrides: { root: { backgroundImage: "none" } } },
      MuiTextField: { defaultProps: { size: "small" } },
    },
  });
}
