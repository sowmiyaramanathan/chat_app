import { Box, Typography } from "@mui/material";
import ContactItem from "./ContactItem";
import { UserSummary } from "./types";
import { STRINGS } from "./keys";
import { panelHeader } from "./styles";

export default function ContactList({
  contacts,
  selectedContactId,
  onContactSelect,
}: {
  contacts: UserSummary[];
  selectedContactId: number | null;
  onContactSelect: (contact: UserSummary, friend: boolean) => void;
}) {
  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Box sx={panelHeader}>
        <Typography variant="subtitle1" color="primary.main" sx={{ fontWeight: 700 }}>
          {STRINGS.chat.contactsTitle}
        </Typography>
      </Box>
      <Box
        sx={{
          flex: 1,
          overflow: "auto",
          "&::-webkit-scrollbar": { width: 6 },
          "&::-webkit-scrollbar-thumb": {
            bgcolor: "secondary.dark",
            borderRadius: 3,
          },
        }}
      >
        {contacts.length === 0 ? (
          <Box sx={{ p: 3, textAlign: "center" }}>
            <Typography variant="body2" color="text.secondary">
              {STRINGS.chat.contactsEmpty}
            </Typography>
          </Box>
        ) : (
          contacts.map((contact) => (
            <ContactItem
              key={contact.ID}
              contact={contact}
              isSelected={selectedContactId === contact.ID}
              onSelect={onContactSelect}
            />
          ))
        )}
      </Box>
    </Box>
  );
}
