import Chats from "../../../../components/Chats";
import { usePagedUsers } from "../../../../components/usePagedUsers";

export default function chats() {
  const { users, error, loading, loadingMore, hasNextPage, loadMore } = usePagedUsers("/user/friends");

  return <Chats contacts={users} contactsError={loading ? null : error} knownFriend hasNextPage={hasNextPage} loadingMore={loadingMore} onLoadMore={loadMore} />;
}
