package simulation

import (
	"fmt"
	"math/big"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stackup-wallet/stackup-bundler/pkg/tracer"
	"github.com/stackup-wallet/stackup-bundler/pkg/userop"
)

type storageSlots mapset.Set[string]

type storageSlotsByEntity map[common.Address]storageSlots

func newStorageSlotsByEntity(stakes EntityStakes, keccak []string) storageSlotsByEntity {
	storageSlotsByEntity := make(storageSlotsByEntity)

	for _, k := range keccak {
		value := common.Bytes2Hex(crypto.Keccak256(common.Hex2Bytes(k[2:])))

		for addr := range stakes {
			if addr == common.HexToAddress("0x") {
				continue
			}
			if _, ok := storageSlotsByEntity[addr]; !ok {
				storageSlotsByEntity[addr] = mapset.NewSet[string]()
			}

			addrPadded := hexutil.Encode(common.LeftPadBytes(addr.Bytes(), 32))
			if strings.HasPrefix(k, addrPadded) {
				storageSlotsByEntity[addr].Add(value)
			}
		}
	}

	return storageSlotsByEntity
}

func isAssociatedWith(slots storageSlots, slot string) bool {
	slotN, _ := big.NewInt(0).SetString(fmt.Sprintf("0x%s", slot), 0)
	for _, k := range slots.ToSlice() {
		kn, _ := big.NewInt(0).SetString(fmt.Sprintf("0x%s", k), 0)
		if slotN.Cmp(kn) >= 0 && slotN.Cmp(big.NewInt(0).Add(kn, big.NewInt(128))) <= 0 {
			return true
		}
	}

	return false
}

func validateStorageSlotsForEntity(
	entityName string,
	op *userop.UserOperation,
	entryPoint common.Address,
	slotsByEntity storageSlotsByEntity,
	entityAccess tracer.AccessMap,
	entityAddr common.Address,
	entityIsStaked bool,
) error {
	senderSlots, senderSlotOk := slotsByEntity[op.Sender]
	if !senderSlotOk {
		senderSlots = mapset.NewSet[string]()
	}
	storageSlots, entitySlotOk := slotsByEntity[entityAddr]
	if !entitySlotOk {
		storageSlots = mapset.NewSet[string]()
	}

	for addr, access := range entityAccess {
		if addr == op.Sender || addr == entryPoint {
			continue
		}

		var mustStakeSlot string
		accessTypes := map[string]tracer.Counts{
			"read":  access.Reads,
			"write": access.Writes,
		}
		for key, slotCount := range accessTypes {
			for slot := range slotCount {
				if isAssociatedWith(senderSlots, slot) {
					// if len(op.InitCode) > 0 {
					// 	mustStakeSlot = slot
					// } else {
					// 	continue
					// }
					continue
				} else if isAssociatedWith(storageSlots, slot) || addr == entityAddr {
					mustStakeSlot = slot
				} else {
					return fmt.Errorf("%s has forbidden %s to %s slot %s", entityName, key, addr2KnownEntity(op, addr), slot)
				}
			}
		}

		if mustStakeSlot != "" && !entityIsStaked {
			return fmt.Errorf(
				"unstaked %s accessed %s slot %s",
				entityName,
				addr2KnownEntity(op, addr),
				mustStakeSlot,
			)
		}
	}

	return nil
}
