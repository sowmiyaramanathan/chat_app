import { Box, Button, Typography } from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import { setPvtKey, setToken } from "../token/token";
import { STRINGS } from "./keys";
import { containedButton, panelCard } from "./styles";

function signout() {
  setToken(null);
  setPvtKey(null);
  window.location.reload();
}

export default function Profile({ name }: { name: string }) {
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
