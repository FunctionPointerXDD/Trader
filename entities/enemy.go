package entities

import "github.com/FunctionPointerXDD/Trader/components"

type Enemy struct {
	*Sprite       // struct embeding -> Enemy를 Sprite처럼 사용가능.
	FollowsPlayer bool
	CombatComp    *components.EnemyCombat
}
