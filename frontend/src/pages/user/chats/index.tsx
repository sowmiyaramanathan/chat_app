import { useState, useEffect } from "react";
import axios from "axios";
import Chats from "../../../../components/Chats";
import { UserSummary } from "../../../../components/types";
import { getApiErrorMessage } from "../../../../components/api";
import { STRINGS } from "../../../../components/keys";
import { API_BASE_URL } from "../../../../components/config";

export default function chats() {
  const [contacts, setContacts] = useState<UserSummary[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token =
      typeof window !== "undefined" && localStorage.getItem("token");
    axios
      .get(`${API_BASE_URL}/user/users`, {
        headers: {
          Authorization: "Bearer " + token,
        },
      })
      .then((response) => {
        setContacts(response.data);
      })
      .catch((requestError: unknown) => {
        setError(getApiErrorMessage(requestError, STRINGS.errors.loadContacts));
      });
  }, []);

  return (
    <>
      <Chats contacts={contacts} contactsError={error} />
    </>
  );
}
