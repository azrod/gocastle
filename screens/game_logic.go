package screens

import (
	"fmt"
	"gocastle/model"
	"math/rand"

	"fyne.io/fyne/v2/dialog"
)

// actOnDirectionKey take player's new coordinates and act on it
func actOnDirectionKey(newX, newY int) {
	// before doing anything, check if we aren't out of bounds
	if currentMap.CheckOutOfBounds(newX, newY) {
		// Player tries to escape map, prevent this, lose 2 seconds
		addLogEntry("you are blocked!")
		model.IncrementTimeSinceBegin(2)
	} else {
		// let's check if we find a NPC on our path
		if npcId := currentMap.NPCList.GetNPCAtPosition(newX, newY); npcId != -1 {
			// get the real NPC from the list, not a copy
			// TODO improve this
			npc := &currentMap.NPCList.List[npcId]
			// if yes, is the NPC hostile?
			if npc.Hostile {
				// let's attack!
				// TODO make this depending on gear
				addLogEntry(npc.HandleNPCDamage(player.PhysicalDamage))
				npc.CurrentHP = npc.CurrentHP - player.PhysicalDamage
				if npc.IsNPCDead() {
					if player.ChangeXP(npc.LootXP) {
						levelUpEntry := fmt.Sprintf("Level up! You are now level %d", player.Level)
						addLogEntry(levelUpEntry)
						levelUpPopup := showLevelUpScreen()
						dialog.ShowCustomConfirm("Level up!", "Validate", "Close", levelUpPopup, func(validate bool) {
							player.RefreshStats(true)
							updateStatsArea()
						}, currentWindow)
					}
					player.ChangeGold(npc.LootGold)
					currentMap.NPCList.RemoveNPCByIndex(npcId)
				}
				// attacking costs 5 seconds
				model.IncrementTimeSinceBegin(5)

			} else {
				// NPC is not hostile, we don't want to hurt them, but lost 2s
				addLogEntry("you are blocked!")
				model.IncrementTimeSinceBegin(2)
			}
		} else {
			// no NPC found on our path, let's check if we can move
			if currentMap.CheckTileIsWalkable(newX, newY) {
				// path is free, let's move (3sec cost)
				player.Avatar.MoveAvatar(newX, newY)
				model.IncrementTimeSinceBegin(3)
			} else {
				// you "hit" a wall, but lost 2s
				addLogEntry("you are blocked!")
				model.IncrementTimeSinceBegin(2)
			}
		}
	}
}

// newTurnForNPCs manages all the map's NPCs actions
func newTurnForNPCs() {
	// for all NPCs, move on a random adjacent tile
	for index := range currentMap.NPCList.List {
		npc := &currentMap.NPCList.List[index]

		var newX, newY int
		if npc.Hostile && npc.Avatar.DistanceFromAvatar(&player.Avatar) <= 10 {
			// player is near, move toward him/her
			newX, newY = npc.Avatar.MoveAvatarTowardsAvatar(&player.Avatar)
		} else {
			// move randomly
			newX = npc.Avatar.PosX + rand.Intn(3) - 1
			newY = npc.Avatar.PosY + rand.Intn(3) - 1

		}

		// don't check / try to move if coordinates stay the same
		if newX != npc.Avatar.PosX || newY != npc.Avatar.PosY {
			// before doing anything, check if we aren't out of bounds
			if !currentMap.CheckOutOfBounds(newX, newY) {
				// let's check if we find another NPC on our NPC's path
				if npcId := currentMap.NPCList.GetNPCAtPosition(newX, newY); npcId != -1 {
					otherNPC := &currentMap.NPCList.List[npcId]
					if (npc.Hostile && !otherNPC.Hostile) ||
						(!npc.Hostile && otherNPC.Hostile) {
						// TODO hostile NPC should attack friendly NPC
						// and vice versa
						addLogEntry(fmt.Sprintf("%s tries to attack %s", npc.Name, otherNPC.Name))
					}
					// let's then check we don't collide with player
				} else if player.Avatar.CollideWithPlayer(newX, newY) {
					if npc.Hostile {
						// TODO hostile NPC should attack player
						addLogEntry(fmt.Sprintf("%s tries to attack you", npc.Name))
					}
					// no ones in our NPC's way
				} else if currentMap.CheckTileIsWalkable(newX, newY) {
					npc.Avatar.MoveAvatar(newX, newY)
				}
			}
		}
	}
}