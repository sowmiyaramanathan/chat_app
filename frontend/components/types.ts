export interface UserSummary {
  ID: number;
  Username: string;
}

export interface ChatMessage {
  ID: number;
  FromUserID: number;
  ToUserID: number;
  Message: string;
  CreatedAt?: string;
}

export interface ChatPageInfo {
  HasNextPage: boolean;
  EndCursor: string;
}

export interface ChatMessagesResponse {
  Messages: ChatMessage[];
  PageInfo: ChatPageInfo;
}

export interface FriendRequest {
  FromUserID: number;
  Username: string;
}
