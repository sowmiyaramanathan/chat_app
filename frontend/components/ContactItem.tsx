import { Box, Typography } from "@mui/material";
import axios from "axios";
import { UserSummary } from "./types";

export default function ContactItem({
  contact,
  onSelect,
}: {
  contact: UserSummary;
  onSelect: (contact: UserSummary, friend: boolean) => void;
}) {
  const handleClick = () => {
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
        if (resp.data.data === true) {
          onSelect(contact, true);
        } else {
          onSelect(contact, false);
        }
      })
      .catch((err) => {
        console.log(err);
      });
  };

  return (
    <Box
      sx={{
        p: 2,
        borderBottom: "1px solid",
        borderColor: "msgBg.main",
        cursor: "pointer",
      }}
      onClick={handleClick}
    >
      <Typography variant="h6" color="primary.light">
        {contact.Username}
      </Typography>
    </Box>
  );
}
