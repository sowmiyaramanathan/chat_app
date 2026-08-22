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

export const theme = createTheme({
  palette: {
    primary: {
      main: "#5C7A5E",
      light: "#FFF8E7",
      dark: "#3D5240",
      contrastText: "#FFF8E7",
    },
    secondary: {
      main: "#C9E4CA",
      light: "#E8F4EA",
      dark: "#8FB892",
      contrastText: "#3D5240",
    },
    msgBg: {
      main: "#A8D8EA",
      light: "#D4ECF7",
      dark: "#7BB8D4",
      contrastText: "#2C3E50",
    },
    surface: {
      main: "#FFFCF5",
      light: "#FFFFFF",
      dark: "#F5F0E8",
      contrastText: "#4A4A3A",
    },
    accent: {
      main: "#E8A87C",
      light: "#F4C4A8",
      dark: "#C8865A",
      contrastText: "#4A4A3A",
    },
    background: {
      default: "#E8F4EA",
      paper: "#FFFCF5",
    },
    text: {
      primary: "#4A4A3A",
      secondary: "#6B6B5B",
    },
    error: {
      main: "#D4726A",
    },
    success: {
      main: "#6B8F71",
    },
  },
  typography: {
    fontFamily: '"Nunito", "Segoe UI", sans-serif',
    h4: { fontWeight: 700, letterSpacing: "-0.02em" },
    h5: { fontWeight: 700 },
    h6: { fontWeight: 600 },
    button: { fontWeight: 600, textTransform: "none" as const },
  },
  shape: {
    borderRadius: 16,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          padding: "8px 20px",
          boxShadow: "none",
          "&:hover": {
            boxShadow: "0 2px 8px rgba(92, 122, 94, 0.2)",
          },
        },
        contained: {
          "&:hover": {
            backgroundColor: "#4A6A4C",
          },
        },
      },
    },
    MuiTab: {
      styleOverrides: {
        root: {
          fontWeight: 600,
          fontSize: "0.95rem",
        },
      },
    },
  },
});
