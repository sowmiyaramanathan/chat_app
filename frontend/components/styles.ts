import type { SxProps, Theme } from "@mui/material";

export const panelCard: SxProps<Theme> = {
  bgcolor: "surface.main",
  borderRadius: 2,
  boxShadow: "0 10px 0 rgba(24, 50, 75, 0.08), 0 20px 40px rgba(24, 50, 75, 0.08)",
  border: "1px solid",
  borderColor: "divider",
  overflow: "hidden",
};

export const pageContainer: SxProps<Theme> = {
  maxWidth: { xs: "100%", md: "1100px", lg: "1200px" },
  mx: "auto",
  px: { xs: 2, sm: 3 },
  pt: { xs: 2, md: 4 },
  pb: 4,
};

export const emptyState: SxProps<Theme> = {
  display: "flex",
  flexDirection: "column",
  alignItems: "center",
  justifyContent: "center",
  textAlign: "center",
  gap: 1.25,
  p: 4,
  height: "100%",
  color: "text.secondary",
};

export const outlinedButton: SxProps<Theme> = {
  color: "primary.main",
  borderColor: "primary.main",
  bgcolor: "transparent",
  "&:hover": {
    bgcolor: "secondary.light",
    borderColor: "primary.dark",
  },
};

export const containedButton: SxProps<Theme> = {
  bgcolor: "primary.main",
  color: "primary.contrastText",
  "&:hover": {
    bgcolor: "primary.dark",
  },
};

export const panelHeader: SxProps<Theme> = {
  px: 2.5,
  py: 2,
  bgcolor: "surface.dark",
  borderBottom: "2px solid",
  borderColor: "divider",
};
