import { Box, Stack, Typography } from "@mui/material";
import { useState } from "react";
import ContactList from "./ContactList";
import ChatScreen from "./ChatScreen";
import RequestStatus from "./RequestStatus";
import { UserSummary } from "./types";

export default function Chats({ contacts }: { contacts: UserSummary[] }) {
  const [selectedContact, setSelectedContact] = useState<UserSummary | null>(
    null
  );
  const [isFriend, setIsFriend] = useState<boolean>(false);

  const handleContactSelect = (contact: UserSummary, friend: boolean) => {
    setSelectedContact(contact);
    setIsFriend(friend);
  };

  return (
    <Box sx={{ maxWidth: "80%", mx: "auto", pt: "10vh" }}>
      <Stack direction="row" sx={{ height: "80vh" }}>
        <Box
          sx={{
            flex: 1,
            minWidth: "250px",
            display: "flex",
            flexDirection: "column",
            border: "1px solid",
            borderColor: "msgBg.main",
          }}
        >
          <ContactList
            contacts={contacts}
            onContactSelect={handleContactSelect}
          />
        </Box>
        <Box
          sx={{
            flex: 3,
            borderTop: "1px solid",
            borderRight: "1px solid",
            borderBottom: "1px solid",
            borderColor: "msgBg.main",
          }}
        >
          {selectedContact && isFriend ? (
            <ChatScreen
              toID={selectedContact.ID}
              username={selectedContact.Username}
            />
          ) : selectedContact && !isFriend ? (
            <RequestStatus
              contact={selectedContact}
              onAccept={handleContactSelect}
            />
          ) : (
            <Box
              sx={{
                display: "flex",
                height: "100%",
                flexDirection: "column",
                alignItems: "center",
                justifyContent: "center",
                color: "primary.light",
              }}
            >
              <Typography variant="h6">Start making friends</Typography>
              <Typography variant="body1">
                Pick a contact on your left to start
              </Typography>
            </Box>
          )}
        </Box>
      </Stack>
    </Box>
  );
}
