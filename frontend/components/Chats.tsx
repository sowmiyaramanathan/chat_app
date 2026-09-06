import { Box, Stack, Typography } from "@mui/material";
import ChatBubbleOutlineOutlinedIcon from "@mui/icons-material/ChatBubbleOutlineOutlined";
import { useState } from "react";
import ContactList from "./ContactList";
import ChatScreen from "./ChatScreen";
import RequestStatus from "./RequestStatus";
import { UserSummary } from "./types";
import { STRINGS } from "./keys";
import { emptyState, pageContainer, panelCard } from "./styles";

export default function Chats({ contacts, contactsError }: { contacts: UserSummary[]; contactsError?: string | null }) {
  const [selectedContact, setSelectedContact] = useState<UserSummary | null>(
    null
  );
  const [isFriend, setIsFriend] = useState<boolean>(false);

  const handleContactSelect = (contact: UserSummary, friend: boolean) => {
    setSelectedContact(contact);
    setIsFriend(friend);
  };

  return (
    <Box sx={pageContainer}>
      <Stack
        direction={{ xs: "column", md: "row" }}
        spacing={2}
        sx={{ height: { xs: "auto", md: "78vh" }, minHeight: { md: 520 }, alignItems: "stretch" }}
      >
        <Box
          sx={{
            ...panelCard,
            flex: { xs: "none", md: "0 0 280px" },
            display: "flex",
            flexDirection: "column",
            minHeight: 0,
            height: { xs: 280, md: "100%" },
          }}
        >
          <ContactList
            contacts={contacts}
            selectedContactId={selectedContact?.ID ?? null}
            onContactSelect={handleContactSelect}
            error={contactsError}
          />
        </Box>

        <Box
          sx={{
            ...panelCard,
            flex: 1,
            display: "flex",
            flexDirection: "column",
            minHeight: { xs: 420, md: "100%" },
          }}
        >
          {selectedContact && isFriend ? (
            <ChatScreen
              toID={selectedContact.ID}
              username={selectedContact.UserName}
            />
          ) : selectedContact && !isFriend ? (
            <RequestStatus
              contact={selectedContact}
              onAccept={handleContactSelect}
            />
          ) : (
            <Box sx={emptyState}>
              <ChatBubbleOutlineOutlinedIcon
                sx={{ fontSize: 48, color: "accent.main", mb: 1 }}
              />
              <Typography variant="h6" color="primary.main">
                {STRINGS.chat.emptyTitle}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ maxWidth: 320 }}>
                {STRINGS.chat.emptySubtitle}
              </Typography>
            </Box>
          )}
        </Box>
      </Stack>
    </Box>
  );
}
