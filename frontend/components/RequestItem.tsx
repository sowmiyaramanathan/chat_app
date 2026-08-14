import { Box, Typography, Button } from "@mui/material";
import { FriendRequest } from "./types";

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
    <Box sx={{ display: "flex", p: 2, justifyContent: "space-around" }}>
      <Typography variant="h6" color="primary.light">
        {request.Username}
      </Typography>
      <Box sx={{ display: "flex", gap: 2 }}>
        <Button
          variant="outlined"
          sx={{
            color: "primary.light",
            borderColor: "msgBg.main",
            ":hover": {
              color: "primary.contrastText",
              borderColor: "primary.light",
            },
          }}
          onClick={() => {
            onReject(request.FromUserID);
          }}
        >
          Reject
        </Button>
        <Button
          variant="contained"
          sx={{
            color: "primary.contrastText",
            backgroundColor: "primary.main",
            ":hover": {
              color: "primary.light",
            },
          }}
          onClick={() => {
            onAccept(request.FromUserID);
          }}
        >
          Accept
        </Button>
      </Box>
    </Box>
  );
}
