package entities

type Enemy struct {
	*Sprite       // struct embeding -> Enemy를 Sprite처럼 사용가능.
	FollowsPlayer bool
}
