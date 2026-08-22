import { Box, Button, Typography } from "@mui/material";
import { useRouter } from "next/router";
import { STRINGS } from "./keys";
import { containedButton, panelCard } from "./styles";

export default function Home() {
  const router = useRouter();

  function signup() {
    router.push({ pathname: "/user/signup" });
  }

  return (
    <Box
      sx={{
        display: "flex",
        minHeight: { xs: "70vh", md: "75vh" },
        justifyContent: "center",
        alignItems: "center",
        px: 2,
      }}
    >
      <Box
        sx={{
          ...panelCard,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          textAlign: "center",
          gap: 2,
          p: { xs: 4, md: 6 },
          maxWidth: 480,
        }}
      >
        <Typography
          variant="h4"
          color="primary.main"
          sx={{ fontWeight: 700 }}
        >
          {STRINGS.home.heading}
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ maxWidth: 360 }}>
          {STRINGS.home.tagline}
        </Typography>
        <Button variant="contained" onClick={signup} sx={containedButton}>
          {STRINGS.home.cta}
        </Button>
      </Box>
    </Box>
  );
}
