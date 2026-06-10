package main

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

const (
	maskHasGun    uint32 = 0x0001 // 0000 0000 0000 0000 0000 0000 0000 0001
	maskHasHouse  uint32 = 0x0002 // 0000 0000 0000 0000 0000 0000 0000 0010
	maskHasFamily uint32 = 0x0004 // 0000 0000 0000 0000 0000 0000 0000 0100

	maskGamePersonType  uint32 = 0x0018 // 0000 0000 0000 0000 0000 0000 0001 1000
	shiftGamePersonType        = 3

	maskManaHealth uint32 = 0x03FF // 0000 0000 0000 0000 0000 0011 1111 1111
	shiftHealth           = 5
	shiftMana             = 15

	maskLevel    uint16 = 0x000F // 0000 0000 0000 1111
	maskExp      uint16 = 0x00F0 // 0000 0000 1111 0000
	maskStrength uint16 = 0x0F00 // 0000 1111 0000 0000
	maskRespect  uint16 = 0xF000 // 1111 0000 0000 0000

	shiftLevel    = 0
	shiftExp      = 4
	shiftStrength = 8
	shiftRespect  = 12
)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = (person.manaHealthFlags & ^maskManaHealth) | (uint32(mana) << shiftMana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = (person.manaHealthFlags & ^maskManaHealth) | (uint32(health) << shiftHealth)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.packedStats = (person.packedStats & ^maskRespect) | (uint16(respect) << shiftRespect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.packedStats = (person.packedStats & ^maskStrength) | (uint16(strength) << shiftStrength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.packedStats = (person.packedStats & ^maskExp) | (uint16(experience) << shiftExp)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.packedStats = (person.packedStats & ^maskLevel) | (uint16(level) << shiftLevel)

	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags |= maskHasHouse
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags |= maskHasGun
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags |= maskHasFamily
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = (person.manaHealthFlags & ^maskGamePersonType) | (uint32(personType) << shiftGamePersonType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x    int32  // 4 byte
	y    int32  // 4 byte
	z    int32  // 4 byte
	gold uint32 // 0..2_000_000 = 21 bites // 11 free bites // 4 byte

	manaHealthFlags uint32 // 7 free bites, mana (0..1000 = 10 bites), health (0..1000 = 10 bites), type (0..2 = 2 bites), family (0..1 = 1 bite), house (0..1 = 1 bite), gun (0..1 = 1 bite)
	packedStats     uint16 // 0 free bites, respect (0..10 = 4 bites), strength (0..10 = 4 bites), experience (0..10 = 4 bites), level (0..10 = 4 bites)

	name [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	var person GamePerson
	for _, option := range options {
		option(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int((p.manaHealthFlags >> shiftMana) & 0x03FF)
}

func (p *GamePerson) Health() int {
	return int((p.manaHealthFlags >> shiftHealth) & 0x03FF)
}

func (p *GamePerson) Respect() int {
	return int((p.packedStats >> shiftRespect) & 0x0F)
}

func (p *GamePerson) Strength() int {
	return int((p.packedStats >> shiftStrength) & 0x0F)
}

func (p *GamePerson) Experience() int {
	return int((p.packedStats >> shiftExp) & 0x0F)
}

func (p *GamePerson) Level() int {
	return int((p.packedStats >> shiftLevel) & 0x0F)
}

func (p *GamePerson) HasHouse() bool {
	return uint32(p.manaHealthFlags)&(maskHasHouse) != 0
}

func (p *GamePerson) HasGun() bool {
	return uint32(p.manaHealthFlags)&(maskHasGun) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return uint32(p.manaHealthFlags)&(maskHasFamily) != 0
}

func (p *GamePerson) Type() int {
	return int(((p.manaHealthFlags) >> shiftGamePersonType) & uint32(24))
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	fmt.Printf("SIZE %d\n", unsafe.Sizeof(person))
	fmt.Printf("%+v\n", person)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
