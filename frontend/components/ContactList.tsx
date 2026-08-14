import { Box } from "@mui/material";
import ContactItem from "./ContactItem";
import { UserSummary } from "./types";

export default function ContactList({
  contacts,
  onContactSelect,
}: {
  contacts: UserSummary[];
  onContactSelect: (contact: UserSummary, friend: boolean) => void;
}) {
  return (
    <Box>
      {contacts.map((contact) => (
        <ContactItem
          key={contact.ID}
          contact={contact}
          onSelect={onContactSelect}
        />
      ))}
    </Box>
  );
}
