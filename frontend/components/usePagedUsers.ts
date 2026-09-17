import axios from "axios";
import { useCallback, useEffect, useState } from "react";
import { API_BASE_URL } from "./config";
import { getApiErrorMessage } from "./api";
import { STRINGS } from "./keys";
import { UserPage, UserSummary } from "./types";

export function usePagedUsers(path: "/user/friends" | "/user/users") {
  const [users, setUsers] = useState<UserSummary[]>([]);
  const [cursor, setCursor] = useState("");
  const [hasNextPage, setHasNextPage] = useState(false);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadPage = useCallback(async (nextCursor = "") => {
    const isFirstPage = nextCursor === "";
    isFirstPage ? setLoading(true) : setLoadingMore(true);
    setError(null);
    try {
      const response = await axios.get<UserPage>(`${API_BASE_URL}${path}`, {
        params: { limit: 25, ...(nextCursor ? { cursor: nextCursor } : {}) },
        headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      });
      setUsers((current) => isFirstPage ? response.data.Users : [...current, ...response.data.Users]);
      setCursor(response.data.PageInfo.EndCursor);
      setHasNextPage(response.data.PageInfo.HasNextPage);
    } catch (requestError: unknown) {
      setError(getApiErrorMessage(requestError, STRINGS.errors.loadContacts));
    } finally {
      isFirstPage ? setLoading(false) : setLoadingMore(false);
    }
  }, [path]);

  useEffect(() => {
    void loadPage();
  }, [loadPage]);

  return { users, error, loading, loadingMore, hasNextPage, loadMore: () => void loadPage(cursor) };
}
