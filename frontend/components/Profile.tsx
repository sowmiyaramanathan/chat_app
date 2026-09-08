import { Alert, Box, Button, Typography } from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import { setToken } from "../token/token";
import { STRINGS } from "./keys";
import { containedButton, panelCard } from "./styles";

function signout() {
  setToken(null);
  window.location.reload();
}

export default function Profile({ name, error }: { name: string; error?: string | null }) {
  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        px: 2,
        pt: { xs: 4, md: 6 },
      }}
    >
      <Box
        sx={{
          ...panelCard,
          display: "flex",
          alignItems: "center",
          flexDirection: "column",
          gap: 3,
          p: { xs: 4, md: 5 },
          maxWidth: 420,
          width: "100%",
        }}
      >
        {error && <Alert severity="error" sx={{ width: "100%" }}>{error}</Alert>}
        <Typography variant="h5" color="primary.main" sx={{ textAlign: "center" }}>
          {STRINGS.profile.welcome(name)}
        </Typography>
        <Button
          variant="contained"
          onClick={signout}
          startIcon={<LogoutIcon />}
          sx={containedButton}
        >
          {STRINGS.profile.signOut}
        </Button>
      </Box>
    </Box>
  );
}
