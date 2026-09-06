import { Alert, Box, CircularProgress, Typography } from "@mui/material";
import PersonAddOutlinedIcon from "@mui/icons-material/PersonAddOutlined";
import RequestItem from "./RequestItem";
import { useCallback, useEffect, useState } from "react";
import { acceptRequest, fetchRequests, getApiErrorMessage, rejectRequest } from "./api";
import { FriendRequest } from "./types";
import { STRINGS } from "./keys";
import { emptyState, pageContainer, panelCard, panelHeader } from "./styles";

export default function RequestList() {
  const [requests, setRequests] = useState<FriendRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const refreshRequests = useCallback(() => {
    setLoading(true);
    setError(null);
    fetchRequests()
      .then((res) => setRequests(res.data))
      .catch((err: unknown) => {
        setError(getApiErrorMessage(err, STRINGS.errors.loadRequests));
      })
      .finally(() => setLoading(false));
  }, []);

  const handleAccept = (id: string) => {
    setError(null);
    acceptRequest(id)
      .then(() => refreshRequests())
      .catch((err: unknown) => setError(getApiErrorMessage(err, STRINGS.errors.updateRequest)));
  };

  const handleReject = (id: string) => {
    setError(null);
    rejectRequest(id)
      .then(() => refreshRequests())
      .catch((err: unknown) => setError(getApiErrorMessage(err, STRINGS.errors.updateRequest)));
  };

  useEffect(() => {
    refreshRequests();
  }, [refreshRequests]);

  return (
    <Box sx={pageContainer}>
      <Box
        sx={{
          ...panelCard,
          height: { xs: "auto", md: "78vh" },
          minHeight: { md: 520 },
          display: "flex",
          flexDirection: "column",
        }}
      >
        <Box sx={panelHeader}>
          <Typography variant="h6" color="primary.main" sx={{ fontWeight: 700 }}>
            {STRINGS.requests.pageTitle}
          </Typography>
        </Box>

        {loading ? (
          <Box
            sx={{
              flex: 1,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              gap: 2,
              p: 4,
            }}
          >
            <CircularProgress size={28} sx={{ color: "primary.main" }} />
          </Box>
        ) : error ? (
          <Box sx={{ p: 3 }}>
            <Alert severity="error" sx={{ borderRadius: 2 }}>
              {error}
            </Alert>
          </Box>
        ) : requests.length === 0 ? (
          <Box sx={emptyState}>
            <PersonAddOutlinedIcon
              sx={{ fontSize: 48, color: "accent.main", mb: 1 }}
            />
            <Typography variant="h6" color="primary.main">
              {STRINGS.requests.emptyTitle}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 360 }}>
              {STRINGS.requests.emptySubtitle}
            </Typography>
          </Box>
        ) : (
          <Box sx={{ flex: 1, overflow: "auto" }}>
            {requests.map((request) => (
              <RequestItem
                key={request.FromUserID}
                request={request}
                onAccept={handleAccept}
                onReject={handleReject}
              />
            ))}
          </Box>
        )}
      </Box>
    </Box>
  );
}
