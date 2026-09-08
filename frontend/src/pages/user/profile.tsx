import { useState, useEffect } from "react";
import axios from "axios";
import Profile from "../../../components/Profile";
import { getApiErrorMessage } from "../../../components/api";
import { STRINGS } from "../../../components/keys";
import { API_BASE_URL } from "../../../components/config";
export default function profile() {
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token =
      typeof window !== "undefined" && localStorage.getItem("token");
    axios
      .get(`${API_BASE_URL}/user/profile`, {
        headers: {
          Authorization: "Bearer " + token,
        },
      })
      .then((response) => {
        setName(response.data.name);
      })
      .catch((requestError: unknown) => {
        setError(getApiErrorMessage(requestError, STRINGS.errors.loadProfile));
      });
  }, []);

  return (
    <>
      <Profile name={name} error={error} />
    </>
  );
}
