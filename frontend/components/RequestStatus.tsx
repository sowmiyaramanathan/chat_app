import { Box, Button, Typography } from "@mui/material";
import axios from "axios";
import { useEffect, useState } from "react";
import { acceptRequest, rejectRequest } from "./api";
import { UserSummary } from "./types";

export default function RequestStatus({
  contact,
  onAccept,
}: {
  contact: UserSummary;
  onAccept: (contact: UserSummary, friend: boolean) => void;
}) {
  const [status, setStatus] = useState<"" | "Received" | "Sent">("");

  useEffect(() => {
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
            .catch((err) => {
              console.log(err);
            });
        }
      })
      .catch((err) => {
        console.log(err);
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
      .catch((err) => {
        console.log(err);
      });
  }

  return (
    <Box sx={{ display: "flex", flexDirection: "column", flex: 1, height: "80vh" }}>
      <Box sx={{ minHeight: 30, bgcolor: "msgBg.main" }}>
        <Typography variant="h6" color="primary.light" sx={{ p: 1 }}>
          {contact.Username}
        </Typography>
      </Box>
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          flex: 1,
          gap: 3,
        }}
      >
        {status === "Received" ? (
          <>
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
                acceptRequest(contact.ID);
                onAccept(contact, true);
              }}
            >
              Accept
            </Button>
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
                rejectRequest(contact.ID);
                setStatus("");
              }}
            >
              Reject
            </Button>
          </>
        ) : status === "Sent" ? (
          <Typography
            sx={{
              border: "1px solid",
              borderColor: "msgBg.main",
              borderRadius: 1,
              color: "primary.light",
              px: 2,
              py: 0.75,
              fontSize: "0.875rem",
              fontWeight: 500,
              textTransform: "uppercase",
            }}
          >
            Request Pending
          </Typography>
        ) : (
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
            onClick={sendRequest}
          >
            Send Request
          </Button>
        )}
      </Box>
    </Box>
  );
}
