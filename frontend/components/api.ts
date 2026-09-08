import axios from "axios";
import { API_BASE_URL } from "./config";

export interface ApiErrorResponse {
  error?: string;
}

export function getApiErrorCode(error: unknown): string | undefined {
  if (!axios.isAxiosError<ApiErrorResponse>(error)) return undefined;
  return error.response?.data?.error;
}

export function getApiErrorMessage(error: unknown, fallback: string): string {
  const code = getApiErrorCode(error);
  if (code === "invalid token") return "Your session has expired. Please sign in again.";
  if (code === "invalid query params") return "Please check the request and try again.";
  if (code === "username") return "That username could not be found.";
  if (code === "password") return "That password is incorrect.";
  if (code === "number") return "That mobile number is already registered.";
  return fallback;
}

export function acceptRequest(id: string) {
	return axios.put(
		`${API_BASE_URL}/friends/acceptRequest?userID=${id}`,
		null,
		{ headers: { Authorization: `Bearer ${localStorage.getItem("token")}` } }
	);
}

export function rejectRequest(id: string) {
	return axios.put(
		`${API_BASE_URL}/friends/rejectRequest?userID=${id}`,
		null,
		{ headers: { Authorization: `Bearer ${localStorage.getItem("token")}` } }
	);
}

export const fetchRequests = () => {
	return axios.get(`${API_BASE_URL}/friends/getFriendRequests`, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
    },
  });
};
