export const setToken = (token: string | null) => {
  if (token) {
    localStorage.setItem("token", token);
  } else {
    localStorage.removeItem("token");
  }
};

interface TokenPayload {
  userID?: string | number;
  sub?: string | number;
  exp?: number;
}

export const decodeToken = (token: string): TokenPayload | null => {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;

    const normalized = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    const payload = normalized.padEnd(
      normalized.length + ((4 - (normalized.length % 4)) % 4),
      "="
    );
    return JSON.parse(atob(payload)) as TokenPayload;
  } catch {
    return null;
  }
};

export const getTokenUserID = (token: string): string | null => {
  const payload = decodeToken(token);
  const userID = payload?.userID ?? payload?.sub;
  return userID === undefined || userID === null ? null : String(userID);
};

export const getTokenExpiry = (token: string): number | null => {
  const expiry = decodeToken(token)?.exp;
  return typeof expiry === "number" ? expiry : null;
};
