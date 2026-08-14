import { useState, useEffect } from "react";
import axios from "axios";
import Chats from "../../../../components/Chats";
import { UserSummary } from "../../../../components/types";

export default function chats() {
  const [contacts, setContacts] = useState<UserSummary[]>([]);

  useEffect(() => {
    const token =
      typeof window !== "undefined" && localStorage.getItem("token");
    axios
      .get("http://localhost:8000/user/users", {
        headers: {
          Authorization: "Bearer " + token,
        },
      })
      .then((response) => {
        setContacts(response.data);
      })
      .catch((error) => {
        console.log(error.response);
      });
  }, []);

  return (
    <>
      <Chats contacts={contacts} />
    </>
  );
}
