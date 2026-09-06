import {
  Alert,
  Box,
  CircularProgress,
  IconButton,
  Typography,
} from "@mui/material";
import SendIcon from "@mui/icons-material/Send";
import MoodRoundedIcon from "@mui/icons-material/MoodRounded";
import { useEffect, useState, useRef } from "react";
import axios from "axios";
import { MessageField } from "./CustomComponets";
import { ChatMessage, ChatMessagesResponse } from "./types";
import { STRINGS } from "./keys";
import { panelHeader } from "./styles";
import { getApiErrorMessage } from "./api";
import { getTokenUserID } from "../token/token";

function ChatScreen({ toID, username }: { toID: string; username: string }) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [newMessage, setNewMessage] = useState("");
  const [loading, setLoading] = useState(true);
  const [loadingOlder, setLoadingOlder] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [sendError, setSendError] = useState<string | null>(null);
  const messagesContainerRef = useRef<HTMLDivElement | null>(null);
  const ws = useRef<WebSocket | null>(null);
  const currentUserIDRef = useRef<string | null>(null);
  const cursorRef = useRef("");
  const shouldScrollToBottomRef = useRef(false);
  const scrollImmediatelyRef = useRef(false);
  const historyLoadedRef = useRef(false);
  const pendingRealtimeMessagesRef = useRef<ChatMessage[]>([]);
  const pageSize = 50;

  const normalizeMessage = (message: ChatMessage): ChatMessage => ({
    ...message,
    ID: Number(message.ID ?? 0),
    FromUserID: message.FromUserID,
    ToUserID: message.ToUserID,
  });

  const appendRealtimeMessage = (message: ChatMessage) => {
    const currentUserID = currentUserIDRef.current;
    const belongsToConversation =
      currentUserID !== null &&
      ((message.FromUserID === currentUserID && message.ToUserID === toID) ||
        (message.FromUserID === toID && message.ToUserID === currentUserID));

    if (!belongsToConversation) return;

    const container = messagesContainerRef.current;
    const isNearBottom = container
      ? container.scrollHeight - container.scrollTop - container.clientHeight < 80
      : true;
    shouldScrollToBottomRef.current = isNearBottom;
    const normalizedMessage = normalizeMessage(message);
    if (!historyLoadedRef.current) {
      pendingRealtimeMessagesRef.current.push(normalizedMessage);
      return;
    }
    setMessages((previous) => [...previous, normalizedMessage]);
  };

  useEffect(() => {
    const token =
      typeof window !== "undefined" ? localStorage.getItem("token") : null;
    const currentUserID = token ? getTokenUserID(token) : null;
    currentUserIDRef.current = currentUserID;

    setLoading(true);
    setError(null);

    if (!token || !currentUserID) {
      setError(STRINGS.errors.invalidToken);
      setLoading(false);
      return;
    }

    let cancelled = false;
    cursorRef.current = "";
    historyLoadedRef.current = false;
    pendingRealtimeMessagesRef.current = [];
    setMessages([]);
    setHasMore(false);

    axios
      .get<ChatMessagesResponse>("http://localhost:8000/message/view", {
        params: { toID: toID, limit: pageSize },
        headers: {
          Authorization: token ? `Bearer ${token}` : "",
        },
      })
      .then((response) => {
        if (cancelled) return;
        const page = response.data;
        const pageMessages = (page.Messages ?? []).map(normalizeMessage);
        const pendingMessages = pendingRealtimeMessagesRef.current;
        pendingRealtimeMessagesRef.current = [];
        setMessages([...pageMessages.reverse(), ...pendingMessages]);
        historyLoadedRef.current = true;
        cursorRef.current = page.PageInfo?.EndCursor ?? "";
        setHasMore(page.PageInfo?.HasNextPage ?? false);
        shouldScrollToBottomRef.current = true;
        scrollImmediatelyRef.current = true;
      })
      .catch(() => {
        if (cancelled) return;
        setError(STRINGS.errors.loadMessages);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });

    ws.current = new WebSocket(
      `ws://localhost:8000/ws/${currentUserID}?token=${token}`
    );

    ws.current.onmessage = (event) => {
      for (const frame of String(event.data).split("\n")) {
        try {
          appendRealtimeMessage(JSON.parse(frame) as ChatMessage);
        } catch {
          // Ignore malformed frames without breaking the websocket listener.
        }
      }
    };

    return () => {
      cancelled = true;
      ws.current?.close();
      ws.current = null;
    };
  }, [toID]);

  useEffect(() => {
    if (!shouldScrollToBottomRef.current) return;

    const container = messagesContainerRef.current;
    const scrollImmediately = scrollImmediatelyRef.current;
    shouldScrollToBottomRef.current = false;
    scrollImmediatelyRef.current = false;

    requestAnimationFrame(() => {
      if (!container) return;
      if (scrollImmediately) {
        container.scrollTop = container.scrollHeight;
        return;
      }
      container.scrollTo({ top: container.scrollHeight, behavior: "smooth" });
    });
  }, [messages]);

  const loadOlderMessages = async () => {
    if (loadingOlder || !hasMore || !cursorRef.current) return;

    const container = messagesContainerRef.current;
    const previousHeight = container?.scrollHeight ?? 0;
    const previousTop = container?.scrollTop ?? 0;
    const token = localStorage.getItem("token");
    if (!token) return;

    setLoadingOlder(true);
    try {
      const response = await axios.get<ChatMessagesResponse>(
        "http://localhost:8000/message/view",
        {
          params: {
            toID: toID,
            limit: pageSize,
            cursor: cursorRef.current,
          },
          headers: { Authorization: `Bearer ${token}` },
        }
      );

      const page = response.data;
      const olderMessages = (page.Messages ?? [])
        .map(normalizeMessage)
        .reverse();
      setMessages((previous) => [...olderMessages, ...previous]);
      cursorRef.current = page.PageInfo?.EndCursor ?? "";
      setHasMore(page.PageInfo?.HasNextPage ?? false);

      requestAnimationFrame(() => {
        if (!container) return;
        container.scrollTop = container.scrollHeight - previousHeight + previousTop;
      });
    } catch (requestError: unknown) {
      setError(getApiErrorMessage(requestError, STRINGS.errors.loadMessages));
    } finally {
      setLoadingOlder(false);
    }
  };

  const handleSendMessage = async () => {
    if (newMessage.trim() === "") return;
    setSendError(null);

    const token = localStorage.getItem("token");
    const currentUserID =
      currentUserIDRef.current ?? (token ? getTokenUserID(token) : null);

    if (!token || !currentUserID) {
      setError(STRINGS.errors.invalidToken);
      return;
    }

    try {
      await axios.post(
        `http://localhost:8000/message/create?toID=${toID}`,
        { message: newMessage },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setNewMessage("");
    } catch (requestError: unknown) {
      setSendError(getApiErrorMessage(requestError, STRINGS.errors.sendMessage));
    }
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", flex: 1, height: "100%" }}>
      <Box sx={{ ...panelHeader, display: "flex", alignItems: "center", gap: 1.25 }}>
        <Box sx={{ position: "relative", display: "grid", placeItems: "center", width: 42, height: 42, borderRadius: "38% 62% 55% 45% / 52% 42% 58% 48%", bgcolor: "secondary.light", color: "primary.main", transform: "rotate(4deg)" }}>
          <MoodRoundedIcon />
          <Box sx={{ position: "absolute", width: 9, height: 9, right: -1, bottom: 1, borderRadius: "50%", bgcolor: "success.main", border: "2px solid", borderColor: "surface.main" }} />
        </Box>
        <Box sx={{ minWidth: 0 }}>
          <Typography variant="h6" color="text.primary" sx={{ fontWeight: 700, lineHeight: 1.1, overflow: "hidden", textOverflow: "ellipsis" }}>
            {username}
          </Typography>
          <Typography variant="caption" color="success.main" sx={{ fontWeight: 700 }}>online-ish</Typography>
        </Box>
      </Box>

      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          flex: 1,
          p: 2,
          overflow: "hidden",
          bgcolor: "surface.dark",
          backgroundImage: (theme) => `radial-gradient(${theme.palette.secondary.main} 1px, transparent 1px)`,
          backgroundSize: "22px 22px",
          backgroundPosition: "3px 4px",
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
            ref={messagesContainerRef}
            onScroll={(event) => {
              if (event.currentTarget.scrollTop <= 24) {
                void loadOlderMessages();
              }
            }}
            sx={{
              flexGrow: 1,
              minHeight: 0,
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
            {loadingOlder && (
              <Box sx={{ display: "flex", justifyContent: "center", py: 1 }}>
                <CircularProgress size={20} sx={{ color: "primary.main" }} />
              </Box>
            )}
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
                      animation: "floatIn 220ms ease both",
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
                      py: 1.1,
                      borderRadius: message.FromUserID === toID ? "6px 18px 18px 18px" : "18px 6px 18px 18px",
                      wordBreak: "break-word",
                      boxShadow: "0 4px 0 rgba(24, 50, 75, 0.08)",
                    }}
                  >
                    {message.Message}
                  </Typography>
                </Box>
              ))
            )}
          </Box>
        )}

        {sendError && <Alert severity="error" sx={{ mb: 1 }}>{sendError}</Alert>}
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
