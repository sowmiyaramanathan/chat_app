import { Box, Typography, IconButton } from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import { useEffect, useState, useRef } from "react";
import axios from "axios";
import { MessageField } from "./CustomComponets";
import { ChatMessage } from "./types";

function ChatScreen({ toID, username }: { toID: number; username: string }) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    const token =
      typeof window !== "undefined" ? localStorage.getItem("token") : null;

    axios
      .get(`http://localhost:8000/message/view?to_id=${toID}`, {
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
        },
      })
      .then((response) => {
        setMessages(response.data);
      })
      .catch((error) => {
        console.log(error.response);
      });

    ws.current = new WebSocket("ws://localhost:8000/ws");

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
      throw new Error("Invalid token");
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
      } catch (error) {
        console.error("Failed to extract user ID:", error);
      }
    }

    setNewMessage("");
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", flex: 1, height: "80vh" }}>
      <Box sx={{ minHeight: 30, bgcolor: "msgBg.main" }}>
        <Typography variant="h6" color="primary.light" sx={{ p: 1 }}>
          {username}
        </Typography>
      </Box>
      <Box sx={{ display: "flex", flexDirection: "column", flex: 1, p: 2, overflow: "auto" }}>
        <Box
          sx={{
            flexGrow: 1,
            overflow: "auto",
            "&::-webkit-scrollbar": {
              display: "none",
            },
          }}
        >
          {messages.map((message, index) => (
            <Box
              key={index}
              sx={{
                display: "flex",
                justifyContent:
                  message.FromUserID === toID ? "flex-start" : "flex-end",
                mb: 1,
              }}
            >
              <Typography
                variant="body1"
                sx={{
                  backgroundColor:
                    message.FromUserID === toID ? "msgBg.main" : "primary.main",
                  color:
                    message.FromUserID === toID
                      ? "black"
                      : "primary.contrastText",
                  p: 1,
                  borderRadius: 3,
                  wordBreak: "break-word",
                }}
              >
                {message.Message}
              </Typography>
            </Box>
          ))}
          <div ref={messagesEndRef} />
        </Box>
        <Box
          component="form"
          sx={{ display: "flex", alignItems: "center", mt: 2 }}
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
            placeholder="Type your message..."
            sx={{ mr: 1 }}
          />
          <IconButton onClick={handleSendMessage} sx={{ color: "msgBg.main" }}>
            <SendIcon />
          </IconButton>
        </Box>
      </Box>
    </Box>
  );
}

export default ChatScreen;
