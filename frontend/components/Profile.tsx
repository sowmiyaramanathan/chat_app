import { Box, Button, Typography } from "@mui/material";
import LogoutIcon from "@mui/icons-material/Logout";
import { setPvtKey, setToken } from "../token/token";

function signout() {
  setToken(null);
  setPvtKey(null);
  window.location.reload();
}

export default function Profile({ name }: { name: string }) {
  return (
    <Box sx={{ display: "flex", alignItems: "center", flexDirection: "column", pt: 5, gap: 3 }}>
      <Typography color="primary.light"> Welcome {name} </Typography>
      <Button variant="contained" onClick={signout} startIcon={<LogoutIcon />}>
        Sign out
      </Button>
    </Box>
  );
}
