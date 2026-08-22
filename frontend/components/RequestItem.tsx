import { Box, Typography, Button } from "@mui/material";
import { FriendRequest } from "./types";
import { STRINGS } from "./keys";
import { containedButton, outlinedButton } from "./styles";

export default function RequestItem({
  request,
  onAccept,
  onReject,
}: {
  request: FriendRequest;
  onAccept: (id: number) => void;
  onReject: (id: number) => void;
}) {
  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: { xs: "column", sm: "row" },
        alignItems: { xs: "flex-start", sm: "center" },
        justifyContent: "space-between",
        gap: 2,
        p: 2.5,
        borderBottom: "1px solid",
        borderColor: "secondary.light",
        "&:last-child": { borderBottom: "none" },
      }}
    >
      <Box sx={{ display: "flex", alignItems: "center", gap: 1.5 }}>
        <Box
          sx={{
            width: 40,
            height: 40,
            borderRadius: "50%",
            bgcolor: "msgBg.light",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <Typography variant="body1" sx={{ fontWeight: 700, color: "primary.main" }}>
            {request.Username.charAt(0).toUpperCase()}
          </Typography>
        </Box>
        <Typography variant="body1" sx={{ fontWeight: 600, color: "text.primary" }}>
          {request.Username}
        </Typography>
      </Box>
      <Box sx={{ display: "flex", gap: 1.5 }}>
        <Button
          variant="outlined"
          sx={outlinedButton}
          onClick={() => onReject(request.FromUserID)}
        >
          {STRINGS.requests.reject}
        </Button>
        <Button
          variant="contained"
          sx={containedButton}
          onClick={() => onAccept(request.FromUserID)}
        >
          {STRINGS.requests.accept}
        </Button>
      </Box>
    </Box>
  );
}
