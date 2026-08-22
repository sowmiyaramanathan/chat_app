import { Box, Button, CircularProgress, Typography } from "@mui/material";
import PersonOutlinedIcon from "@mui/icons-material/PersonOutlined";
import axios from "axios";
import { useEffect, useState } from "react";
import { acceptRequest, rejectRequest } from "./api";
import { UserSummary } from "./types";
import { STRINGS } from "./keys";
import { containedButton, emptyState, outlinedButton, panelHeader } from "./styles";

export default function RequestStatus({
  contact,
  onAccept,
}: {
  contact: UserSummary;
  onAccept: (contact: UserSummary, friend: boolean) => void;
}) {
  const [status, setStatus] = useState<"" | "Received" | "Sent">("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);

    axios
      .get(
        `http://localhost:8000/friends/isRequestReceived?from_user_id=${contact.ID}`,
        {
          headers: {
            Authorization: "Bearer " + `${localStorage.getItem("token")}`,
          },
        }
      )
      .then((resp) => {
        if (resp.data.data === true) {
          setStatus("Received");
          setLoading(false);
        } else {
          axios
            .get(
              `http://localhost:8000/friends/isRequestSent?to_user_id=${contact.ID}`,
              {
                headers: {
                  Authorization: "Bearer " + `${localStorage.getItem("token")}`,
                },
              }
            )
            .then((resp) => {
              if (resp.data.data === true) {
                setStatus("Sent");
              } else {
                setStatus("");
              }
            })
            .catch((err) => console.log(err))
            .finally(() => setLoading(false));
        }
      })
      .catch((err) => {
        console.log(err);
        setLoading(false);
      });
  }, [contact.ID]);

  function sendRequest() {
    axios
      .post(
        `http://localhost:8000/friends/sendFriendRequest?to_user_id=${contact.ID}`,
        "",
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      )
      .then((resp) => {
        if (resp.data.data === "Sent") {
          setStatus("Sent");
        }
      })
      .catch((err) => console.log(err));
  }

  return (
    <Box sx={{ display: "flex", flexDirection: "column", flex: 1, height: "100%" }}>
      <Box sx={panelHeader}>
        <Typography variant="h6" color="primary.main" sx={{ fontWeight: 700 }}>
          {contact.Username}
        </Typography>
      </Box>

      {loading ? (
        <Box
          sx={{
            flex: 1,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <CircularProgress size={28} sx={{ color: "primary.main" }} />
        </Box>
      ) : (
        <Box sx={{ ...emptyState, bgcolor: "surface.dark" }}>
          <PersonOutlinedIcon
            sx={{ fontSize: 48, color: "accent.main", mb: 1 }}
          />
          <Typography variant="h6" color="primary.main">
            {STRINGS.requests.notFriendTitle}
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 320, mb: 2 }}>
            {STRINGS.requests.notFriendSubtitle}
          </Typography>

          {status === "Received" ? (
            <Box sx={{ display: "flex", gap: 2, flexWrap: "wrap", justifyContent: "center" }}>
              <Button
                variant="contained"
                sx={containedButton}
                onClick={() => {
                  acceptRequest(contact.ID);
                  onAccept(contact, true);
                }}
              >
                {STRINGS.requests.accept}
              </Button>
              <Button
                variant="outlined"
                sx={outlinedButton}
                onClick={() => {
                  rejectRequest(contact.ID);
                  setStatus("");
                }}
              >
                {STRINGS.requests.reject}
              </Button>
            </Box>
          ) : status === "Sent" ? (
            <Typography
              sx={{
                border: "1px solid",
                borderColor: "accent.main",
                borderRadius: 2,
                color: "accent.dark",
                bgcolor: "accent.light",
                px: 2.5,
                py: 1,
                fontSize: "0.875rem",
                fontWeight: 600,
              }}
            >
              {STRINGS.requests.requestPending}
            </Typography>
          ) : (
            <Button variant="contained" sx={containedButton} onClick={sendRequest}>
              {STRINGS.requests.sendRequest}
            </Button>
          )}
        </Box>
      )}
    </Box>
  );
}
