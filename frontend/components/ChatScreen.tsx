import {
  Alert,
  Box,
  CircularProgress,
  IconButton,
  Typography,
} from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import { useEffect, useState, useRef } from "react";
import axios from "axios";
import { MessageField } from "./CustomComponets";
import { ChatMessage } from "./types";
import { STRINGS } from "./keys";
import { panelHeader } from "./styles";

function ChatScreen({ toID, username }: { toID: number; username: string }) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    const token =
      typeof window !== "undefined" ? localStorage.getItem("token") : null;

    setLoading(true);
    setError(null);

    axios
      .get(`http://localhost:8000/message/view?to_id=${toID}`, {
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
        },
      })
      .then((response) => {
        setMessages(response.data);
      })
      .catch((err) => {
        console.log(err.response);
        setError(STRINGS.errors.loadMessages);
      })
      .finally(() => {
        setLoading(false);
      });

    ws.current = new WebSocket(`ws://localhost:8000/ws/${toID}`);

    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data) as ChatMessage;
      setMessages((prevMessages) => [...prevMessages, message]);
    };

    return () => {
      ws.current?.close();
    };
  }, [toID]);

  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  const getUserIDFromToken = (token: string) => {
    const parts = token.split(".");
    if (parts.length !== 3) {
      throw new Error(STRINGS.errors.invalidToken);
    }

    const payload = parts[1];
    const base64 = payload.replace(/-/g, "+").replace(/_/g, "/");
    const jsonString = atob(base64);
    const payloadObj = JSON.parse(jsonString);

    return Number(payloadObj.sub ?? payloadObj.userID);
  };

  const handleSendMessage = async () => {
    if (newMessage.trim() === "") return;

    const token = localStorage.getItem("token");

    await axios.post(
      `http://localhost:8000/message/create?to_id=${toID}`,
      { message: newMessage },
      { headers: { Authorization: `Bearer ${token}` } }
    );

    if (token) {
      try {
        const userID = getUserIDFromToken(token);
        if (ws.current) {
          ws.current.send(
            JSON.stringify({
              Message: newMessage,
              FromUserID: userID,
              ToUserID: toID,
            })
          );
        }
      } catch (err) {
        console.error("Failed to extract user ID:", err);
      }
    }

    setNewMessage("");
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", flex: 1, height: "100%" }}>
      <Box sx={panelHeader}>
        <Typography variant="h6" color="primary.main" sx={{ fontWeight: 700 }}>
          {username}
        </Typography>
      </Box>

      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          flex: 1,
          p: 2,
          overflow: "hidden",
          bgcolor: "surface.dark",
        }}
      >
        {loading ? (
          <Box
            sx={{
              flex: 1,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              gap: 2,
            }}
          >
            <CircularProgress size={28} sx={{ color: "primary.main" }} />
            <Typography color="text.secondary">
              {STRINGS.chat.loadingMessages}
            </Typography>
          </Box>
        ) : error ? (
          <Box sx={{ flex: 1, display: "flex", alignItems: "center", p: 2 }}>
            <Alert severity="error" sx={{ width: "100%", borderRadius: 2 }}>
              {error}
            </Alert>
          </Box>
        ) : (
          <Box
            sx={{
              flexGrow: 1,
              overflow: "auto",
              display: "flex",
              flexDirection: "column",
              gap: 1,
              "&::-webkit-scrollbar": { width: 6 },
              "&::-webkit-scrollbar-thumb": {
                bgcolor: "secondary.dark",
                borderRadius: 3,
              },
            }}
          >
            {messages.length === 0 ? (
              <Box
                sx={{
                  flex: 1,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                <Typography variant="body2" color="text.secondary">
                  {STRINGS.chat.emptyMessages}
                </Typography>
              </Box>
            ) : (
              messages.map((message, index) => (
                <Box
                  key={index}
                  sx={{
                    display: "flex",
                    justifyContent:
                      message.FromUserID === toID ? "flex-start" : "flex-end",
                  }}
                >
                  <Typography
                    variant="body1"
                    sx={{
                      maxWidth: "75%",
                      backgroundColor:
                        message.FromUserID === toID
                          ? "msgBg.main"
                          : "primary.main",
                      color:
                        message.FromUserID === toID
                          ? "msgBg.contrastText"
                          : "primary.contrastText",
                      px: 2,
                      py: 1,
                      borderRadius: 3,
                      wordBreak: "break-word",
                      boxShadow: "0 1px 4px rgba(0,0,0,0.06)",
                    }}
                  >
                    {message.Message}
                  </Typography>
                </Box>
              ))
            )}
            <div ref={messagesEndRef} />
          </Box>
        )}

        <Box
          component="form"
          sx={{
            display: "flex",
            alignItems: "center",
            gap: 1,
            mt: 2,
            pt: 2,
            borderTop: "1px solid",
            borderColor: "secondary.light",
          }}
          onSubmit={(e) => {
            e.preventDefault();
            handleSendMessage();
          }}
        >
          <MessageField
            variant="outlined"
            size="small"
            fullWidth
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder={STRINGS.chat.messagePlaceholder}
          />
          <IconButton
            onClick={handleSendMessage}
            aria-label="Send message"
            sx={{
              bgcolor: "primary.main",
              color: "primary.contrastText",
              "&:hover": { bgcolor: "primary.dark" },
              width: 44,
              height: 44,
            }}
          >
            <SendIcon fontSize="small" />
          </IconButton>
        </Box>
      </Box>
    </Box>
  );
}

export default ChatScreen;
