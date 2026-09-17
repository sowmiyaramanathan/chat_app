import Chats from "../../../components/Chats";
import { usePagedUsers } from "../../../components/usePagedUsers";

export default function Discover() {
  const { users, error, loading, loadingMore, hasNextPage, loadMore } = usePagedUsers("/user/users");

  return <Chats contacts={users} contactsError={loading ? null : error} title="Discover people" knownFriend={false} hasNextPage={hasNextPage} loadingMore={loadingMore} onLoadMore={loadMore} />;
}
