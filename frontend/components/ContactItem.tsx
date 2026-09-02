import { Box, CircularProgress, Typography } from "@mui/material";
import axios from "axios";
import { useState } from "react";
import { UserSummary } from "./types";

export default function ContactItem({
  contact,
  isSelected,
  onSelect,
}: {
  contact: UserSummary;
  isSelected: boolean;
  onSelect: (contact: UserSummary, friend: boolean) => void;
}) {
  const [loading, setLoading] = useState(false);

  const handleClick = () => {
    if (loading) return;
    setLoading(true);

    axios
      .get(
        `http://localhost:8000/friends/isFriend?with_user_id=${contact.ID}`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        }
      )
      .then((resp) => {
        onSelect(contact, resp.data.data === true);
      })
      .catch((err) => {
        console.log(err);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  return (
    <Box
      role="button"
      tabIndex={0}
      aria-selected={isSelected}
      sx={{
        px: 2,
        py: 1.5,
        borderBottom: "1px solid",
        borderColor: "rgba(24, 50, 75, 0.08)",
        cursor: loading ? "wait" : "pointer",
        bgcolor: isSelected ? "secondary.light" : "transparent",
        borderLeft: "4px solid",
        borderLeftColor: isSelected ? "accent.main" : "transparent",
        transition: "background-color 160ms ease, border-color 160ms ease, transform 160ms ease",
        "&:hover": {
          bgcolor: isSelected ? "secondary.light" : "surface.dark",
          transform: "translateX(3px)",
        },
      }}
      onClick={handleClick}
      onKeyDown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          handleClick();
        }
      }}
    >
      <Box sx={{ display: "flex", alignItems: "center", gap: 1.5 }}>
        <Box
          sx={{
            width: 36,
            height: 36,
            borderRadius: "38% 62% 55% 45% / 52% 42% 58% 48%",
            bgcolor: isSelected ? "accent.light" : "msgBg.light",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            flexShrink: 0,
          }}
        >
          <Typography
            variant="body2"
            sx={{ fontWeight: 700, color: "primary.main", textTransform: "uppercase" }}
          >
            {contact.Username.charAt(0)}
          </Typography>
        </Box>
        <Typography
          variant="body1"
          color={isSelected ? "primary.main" : "text.primary"}
          sx={{ flex: 1, fontWeight: isSelected ? 700 : 500 }}
        >
          {contact.Username}
        </Typography>
        {loading && <CircularProgress size={18} sx={{ color: "primary.main" }} />}
      </Box>
    </Box>
  );
}
