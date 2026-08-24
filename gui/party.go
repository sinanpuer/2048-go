package main

import "crypto/rand"

// partyRelayHost is the externally hosted relay all players connect to for
// Battle Royale - it works over the open internet (not just LAN), which
// matters because school/guest WiFi networks commonly block direct
// device-to-device connections even when both players are on the same
// network. The relay's game logic lives in relay/ at the repo root.
const partyRelayHost = "https://merge-kingdom-relay.onrender.com"

const roomCodeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // no 0/O/1/I

// newPartyRoomCode generates a short, human-shareable room code.
func newPartyRoomCode() string {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	code := make([]byte, len(b))
	for i, v := range b {
		code[i] = roomCodeChars[int(v)%len(roomCodeChars)]
	}
	return string(code)
}

// partyRoomURL builds the shareable link everyone (including the creator)
// opens in their browser to play.
func partyRoomURL(roomCode string) string {
	return partyRelayHost + "/?room=" + roomCode
}
