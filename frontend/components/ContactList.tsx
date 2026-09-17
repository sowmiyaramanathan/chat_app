import { Alert, Box, Button, CircularProgress, Typography } from "@mui/material";
import ContactItem from "./ContactItem";
import { UserSummary } from "./types";
import { STRINGS } from "./keys";
import { panelHeader } from "./styles";

export default function ContactList({
  contacts,
  selectedContactId,
  onContactSelect,
  error,
  title = STRINGS.chat.contactsTitle,
  knownFriend,
  hasNextPage = false,
  loadingMore = false,
  onLoadMore,
}: {
  contacts: UserSummary[];
  selectedContactId: string | null;
  onContactSelect: (contact: UserSummary, friend: boolean) => void;
  error?: string | null;
  title?: string;
  knownFriend?: boolean;
  hasNextPage?: boolean;
  loadingMore?: boolean;
  onLoadMore?: () => void;
}) {
  return (
    <Box sx={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Box sx={{ ...panelHeader, display: "flex", alignItems: "baseline", justifyContent: "space-between", gap: 1 }}>
        <Typography variant="h6" color="text.primary" sx={{ fontWeight: 700 }}>
          {title}
        </Typography>
        <Typography variant="caption" color="text.secondary" sx={{ fontWeight: 700 }}>
          {contacts.length} {contacts.length === 1 ? "person" : "people"}
        </Typography>
      </Box>
      {error && <Alert severity="error" sx={{ m: 1.5 }}>{error}</Alert>}
      <Box
        sx={{
          flex: 1,
          minHeight: 0,
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
              friend={knownFriend}
            />
          ))
        )}
        {hasNextPage && (
          <Box sx={{ p: 1.5, display: "flex", justifyContent: "center" }}>
            <Button onClick={onLoadMore} disabled={loadingMore} size="small">
              {loadingMore ? <CircularProgress size={18} /> : "Load more"}
            </Button>
          </Box>
        )}
      </Box>
    </Box>
  );
}
