import { Box, Button, Chip, Stack, Typography } from "@mui/material";
import AutoAwesomeRoundedIcon from "@mui/icons-material/AutoAwesomeRounded";
import ForumRoundedIcon from "@mui/icons-material/ForumRounded";
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
        minHeight: { xs: "72vh", md: "76vh" },
        justifyContent: "center",
        alignItems: "center",
        px: 2,
        py: 5,
      }}
    >
      <Box
        sx={{
          ...panelCard,
          position: "relative",
          overflow: "hidden",
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          textAlign: "center",
          gap: 2.5,
          p: { xs: 4, md: 6 },
          maxWidth: 620,
          width: "100%",
          animation: "floatIn 500ms ease both",
          "&::before": {
            content: '""',
            position: "absolute",
            width: 170,
            height: 170,
            border: "3px dashed",
            borderColor: "accent.main",
            borderRadius: "44% 56% 62% 38% / 44% 38% 62% 56%",
            top: -88,
            right: -55,
            transform: "rotate(18deg)",
            opacity: 0.65,
          },
        }}
      >
        <Chip
          icon={<AutoAwesomeRoundedIcon />}
          label="small talk, big feelings"
          sx={{ bgcolor: "secondary.light", color: "primary.main", fontWeight: 700 }}
        />
        <Box sx={{ display: "grid", placeItems: "center", width: 72, height: 72, borderRadius: "34% 66% 54% 46% / 48% 38% 62% 52%", bgcolor: "accent.main", color: "accent.contrastText", transform: "rotate(-6deg)", boxShadow: "6px 6px 0 rgba(24,50,75,0.12)" }}>
          <ForumRoundedIcon sx={{ fontSize: 38 }} />
        </Box>
        <Typography
          variant="h4"
          color="text.primary"
          sx={{ fontWeight: 700, fontSize: { xs: "2.3rem", md: "3.3rem" }, lineHeight: 1 }}
        >
          {STRINGS.home.heading}
        </Typography>
        <Typography variant="body1" color="text.secondary" sx={{ maxWidth: 410, fontSize: "1.1rem" }}>
          {STRINGS.home.tagline}
        </Typography>
        <Stack direction={{ xs: "column", sm: "row" }} spacing={1.5} sx={{ alignItems: "center" }}>
          <Button variant="contained" onClick={signup} sx={containedButton}>
            {STRINGS.home.cta}
          </Button>
          <Typography variant="caption" color="text.secondary">No awkward icebreakers required.</Typography>
        </Stack>
      </Box>
    </Box>
  );
}
