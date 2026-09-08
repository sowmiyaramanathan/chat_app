const configuredBackendURL =
	process.env.NEXT_PUBLIC_BACKEND_URL ?? "http://localhost:8000";

const backendURL = /^https?:\/\//.test(configuredBackendURL)
	? configuredBackendURL
	: `http://${configuredBackendURL}`;

export const API_BASE_URL = backendURL.replace(/\/+$/, "");
export const WS_BASE_URL = API_BASE_URL.replace(/^http/, "ws");
