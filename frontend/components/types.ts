export interface UserSummary {
  ID: string;
  UserName: string;
}

export interface ChatMessage {
  ID: number;
  FromUserID: string;
  ToUserID: string;
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
  FromUserID: string;
  UserName: string;
}
