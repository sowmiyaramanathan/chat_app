import { Alert, Box, Button, CircularProgress, Typography } from "@mui/material";
import PersonOutlinedIcon from "@mui/icons-material/PersonOutlined";
import axios from "axios";
import { useEffect, useState } from "react";
import { acceptRequest, getApiErrorMessage, rejectRequest } from "./api";
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
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    setError(null);

    axios
      .get(
        `http://localhost:8000/friends/isRequestReceived?userID=${contact.ID}`,
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
              `http://localhost:8000/friends/isRequestSent?userID=${contact.ID}`,
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
            .catch((err: unknown) => setError(getApiErrorMessage(err, STRINGS.errors.loadContacts)))
            .finally(() => setLoading(false));
        }
      })
      .catch((err: unknown) => {
        setError(getApiErrorMessage(err, STRINGS.errors.loadContacts));
        setLoading(false);
      });
  }, [contact.ID]);

  function sendRequest() {
    setError(null);
    axios
      .post(
        `http://localhost:8000/friends/sendRequest?userID=${contact.ID}`,
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
      .catch((err: unknown) => setError(getApiErrorMessage(err, STRINGS.errors.updateRequest)));
  }

  return (
    <Box sx={{ display: "flex", flexDirection: "column", flex: 1, height: "100%" }}>
      <Box sx={panelHeader}>
        <Typography variant="h6" color="primary.main" sx={{ fontWeight: 700 }}>
          {contact.UserName}
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
          {error && <Alert severity="error" sx={{ width: "100%", maxWidth: 420 }}>{error}</Alert>}
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
                  setError(null);
                  acceptRequest(contact.ID)
                    .then(() => onAccept(contact, true))
                    .catch((err: unknown) =>
                      setError(getApiErrorMessage(err, STRINGS.errors.updateRequest))
                    );
                }}
              >
                {STRINGS.requests.accept}
              </Button>
              <Button
                variant="outlined"
                sx={outlinedButton}
                onClick={() => {
                  setError(null);
                  rejectRequest(contact.ID)
                    .then(() => setStatus(""))
                    .catch((err: unknown) =>
                      setError(getApiErrorMessage(err, STRINGS.errors.updateRequest))
                    );
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
