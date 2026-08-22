import type { SxProps, Theme } from "@mui/material";

export const panelCard: SxProps<Theme> = {
  bgcolor: "surface.main",
  borderRadius: 3,
  boxShadow: "0 4px 24px rgba(92, 122, 94, 0.12)",
  border: "1px solid",
  borderColor: "secondary.dark",
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
  gap: 1.5,
  p: 4,
  height: "100%",
  color: "text.secondary",
};

export const outlinedButton: SxProps<Theme> = {
  color: "primary.main",
  borderColor: "primary.main",
  bgcolor: "surface.main",
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
  py: 1.75,
  bgcolor: "secondary.light",
  borderBottom: "1px solid",
  borderColor: "secondary.dark",
};
