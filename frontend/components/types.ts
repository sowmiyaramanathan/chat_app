export interface UserSummary {
  ID: number;
  Username: string;
}

export interface ChatMessage {
  FromUserID: number;
  ToUserID: number;
  Message: string;
}

export interface FriendRequest {
  FromUserID: number;
  Username: string;
}
