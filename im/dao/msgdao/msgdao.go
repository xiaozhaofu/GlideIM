package msgdao

var instance MsgDao

func init() {
	instance = impl{
		GroupMsgDao: GroupMsgDaoImpl,
		ChatMsgDao:  ChatMsgDaoImpl,
		CacheDao:    cacheDao{},
		CommonDao:   commonDao{},
	}
}

type GroupMsgDao interface {
	GetGroupMsgSeq(gid int64) (int64, error)
	UpdateGroupMsgSeq(gid int64, seq int64) error
	CreateGroupMsgSeq(gid int64, step int64) error

	GetMessage(mid int64) (*GroupMessage, error)
	GetGroupMessage(gid int64, page int, pageSize int) ([]*GroupMessage, error)
	GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessage, error)

	AddGroupMessage(message *GroupMessage) error
	UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAt int64, lastMsgSeq int64) error
	GetGroupMessageState(gid int64) (*GroupMessageState, error)

	CreateGroupMemberMsgState(gid int64, uid int64) error
	UpdateGroupMemberMsgState(gid int64, uid int64, lastAck int64, lastAckSeq int64) error
	GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgState, error)
}

type ChatMsgDao interface {
	GetChatMessage(mid ...int64) ([]*ChatMessage, error)
	GetChatMessagesBySession(uid1, uid2 int64, page int, pageSize int) ([]*ChatMessage, error)
	GetRecentChatMessages(uid int64, afterTime int64) ([]*ChatMessage, error)
	AddOrUpdateChatMessage(message *ChatMessage) (bool, error)

	GetChatMessageMidAfter(form, to int64, midAfter int64) ([]*ChatMessage, error)
	GetChatMessageMidSpan(from, to int64, midStart, midEnd int64) ([]*ChatMessage, error)

	AddOfflineMessage(uid int64, mid int64) error
	GetOfflineMessage(uid int64) ([]*OfflineMessage, error)
	DelOfflineMessage(uid int64, mid []int64) error
}

type SessionDao interface {
	GetSession(uid1 int64, uid2 int64) (*Session, error)
	CreateSession(uid1 int64, uid2 int64, updateAt int64) (*Session, error)
	UpdateOrInitSession(uid1 int64, uid2 int64, update int64) error
	GetRecentSession(uid int64, updateAfter int64) ([]*Session, error)
}

type CacheDao interface {
	// GetUserMsgSeq 获取用户全当前局消息 Seq
	GetUserMsgSeq(uid int64) (int64, error)
	// GetIncrUserMsgSeq 返回用户全局消息递增seq, 保证递增, 尽量保持连续, 不保证一定连续
	GetIncrUserMsgSeq(uid int64) (int64, error)
}

type CommonDao interface {
	GetMessageID() (int64, error)
}

type MsgDao interface {
	ChatMsgDao
	GroupMsgDao
	CacheDao
	CommonDao
}

type impl struct {
	ChatMsgDao
	GroupMsgDao
	CacheDao
	CommonDao
}

/////////////////

func GetUserMsgSeq(uid int64) (int64, error) {
	return instance.GetUserMsgSeq(uid)
}

func GetIncrUserMsgSeq(uid int64) (int64, error) {
	return instance.GetIncrUserMsgSeq(uid)
}

/////////////////

func GetMessageID() (int64, error) {
	return instance.GetMessageID()
}

/////////////////

func GetGroupMsgSeq(gid int64) (int64, error) {
	return instance.GetGroupMsgSeq(gid)
}
func UpdateGroupMsgSeq(gid int64, seq int64) error {
	return instance.UpdateGroupMsgSeq(gid, seq)
}
func CreateGroupMsgSeq(gid int64, step int64) error {
	return instance.CreateGroupMsgSeq(gid, step)
}

func GetGroupMessage(mid int64) (*GroupMessage, error) {
	return instance.GetMessage(mid)
}
func GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessage, error) {
	return instance.GetGroupMessageSeqAfter(gid, seqAfter)
}
func AddGroupMessage(message *GroupMessage) error {
	return instance.AddGroupMessage(message)
}
func UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAt int64, lastMsgSeq int64) error {
	return instance.UpdateGroupMessageState(gid, lastMID, lastMsgAt, lastMsgSeq)
}
func GetGroupMessageState(gid int64) (*GroupMessageState, error) {
	return instance.GetGroupMessageState(gid)
}
func UpdateGroupMemberMsgState(gid int64, uid int64, lastAck int64, lastAckSeq int64) error {
	return instance.UpdateGroupMemberMsgState(gid, uid, lastAck, lastAckSeq)
}
func GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgState, error) {
	return instance.GetGroupMemberMsgState(gid, uid)
}

///////////////////////////////////////

func GetChatMessage(mid ...int64) ([]*ChatMessage, error) {
	return instance.GetChatMessage(mid...)
}
func AddChatMessage(message *ChatMessage) (bool, error) {
	return instance.AddOrUpdateChatMessage(message)
}
func GetChatMessageMidAfter(from int64, to int64, midAfter int64) ([]*ChatMessage, error) {
	return instance.GetChatMessageMidAfter(from, to, midAfter)
}
func GetChatMessageMidSpan(from, to int64, midStart, midEnd int64) ([]*ChatMessage, error) {
	return instance.GetChatMessageMidSpan(from, to, midStart, midEnd)
}
func AddOfflineMessage(uid int64, mid int64) error {
	return instance.AddOfflineMessage(uid, mid)
}
func GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	return instance.GetOfflineMessage(uid)
}
func DelOfflineMessage(uid int64, mid []int64) error {
	return instance.DelOfflineMessage(uid, mid)
}
