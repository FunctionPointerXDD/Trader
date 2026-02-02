package components

type Combat interface {
	Health() int
	AttackPower() int
	Attacking() bool
	Attack() bool
	Update()
	Damage(amount int)
}

type BasicCombat struct {
	health      int
	attackPower int
	attacking   bool
}

func NewBasicCombat(health, attackPower int) *BasicCombat {
	return &BasicCombat{
		health,
		attackPower,
		false,
	}
}

// AttackPower implements [Combat].
func (b *BasicCombat) AttackPower() int {
	return b.attackPower
}

// Health implements [Combat].
func (b *BasicCombat) Health() int {
	return b.health
}

// Damage implements [Combat].
func (b *BasicCombat) Damage(amount int) {
	b.health -= amount
}

func (b *BasicCombat) Attacking() bool {
	return b.attacking
}

func (b *BasicCombat) Attack() bool {
	b.attacking = true
	return true
}

func (b *BasicCombat) Update() {
}

// 컴파일러 에러 체크 확인용도(빠진 메서드가 있는지 확인)
var _ Combat = (*BasicCombat)(nil)

type EnemyCombat struct {
	*BasicCombat
	attackCooldown  int
	timeSinceAttack int
}

func NewEnemyCombat(health, attackPower, attackCooldown int) *EnemyCombat {
	return &EnemyCombat{
		NewBasicCombat(health, attackPower),
		attackCooldown,
		0,
	}
}

func (e *EnemyCombat) Attack() bool {
	if e.timeSinceAttack >= e.attackCooldown {
		e.attacking = true
		e.timeSinceAttack = 0
		return true
	}
	return false
}

// ebitengine 특성상 초당 60번 더해짐.. (1초에 60프레임 TPS:60)
func (e *EnemyCombat) Update() {
	e.timeSinceAttack += 1
}

// 컴파일러 에러 체크 확인용도(빠진 메서드가 있는지 확인)
var _ Combat = (*EnemyCombat)(nil)
