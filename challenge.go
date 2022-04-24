package main

import (
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
	"strconv"
)

type ChallengeApiResponse struct {
	CurrentChallenge string `json:"challenge_current"`
	Reward           int    `json:"challenge_reward"`
	Difficulty       int    `json:"challenge_difficulty"`
}

var currentChallenge string = ""

func RefreshChallenge() {
	currentChallenge = strconv.Itoa(int(rand.Uint64()))
}

func GetChallenge() string {
	return currentChallenge
}

func CheckChallenge(magic string, difficulty int) bool {
	hash := sha512.Sum512([]byte(currentChallenge + magic))

	stringhash := hex.EncodeToString(hash[:])
	hexarray := []rune(stringhash)

	for i := 0; i < difficulty; i++ {
		if hexarray[i] != '0' {
			return false
		}
	}

	return true
}
