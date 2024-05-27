package ruleMatch

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"graduationDesign/parser/parserRuler"
	"graduationDesign/structure/ruleStructure"
	"net"
)

// ExecuteLogic evaluates the boolean expression given bools and operators
// in the order of logical AND first, then logical OR.
func ExecuteLogic(bools []bool, operators []int8) (bool, error) {
	if len(bools) == 0 {
		return false, errors.New("the bool array is empty")
	}
	if len(operators) != len(bools)-1 {
		return false, errors.New("the number of operators must be one less than the number of bools")
	}

	// Step 1: Process all && operators
	andProcessed := []bool{bools[0]}
	andOps := []int8{}
	for i := 0; i < len(operators); i++ {
		if operators[i] == 0 {
			andProcessed[len(andProcessed)-1] = andProcessed[len(andProcessed)-1] && bools[i+1]
		} else {
			andProcessed = append(andProcessed, bools[i+1])
			andOps = append(andOps, operators[i])
		}
	}

	// Step 2: Process all || operators
	result := andProcessed[0]
	for i := 0; i < len(andOps); i++ {
		if andOps[i] == 1 {
			result = result || andProcessed[i+1]
		} else {
			return false, fmt.Errorf("invalid operator after && processing: %s", andOps[i])
		}
	}

	return result, nil
}

// 比较两个字节数组是否相等
func bytesEqual(a, b net.HardwareAddr) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var MapLayerTypeToMatchFunc = map[gopacket.LayerType]func(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool{
	layers.LayerTypeEthernet: MACPayloadMatch,
}

func MACPayloadMatch(protocol *ruleStructure.Protocol, payloadPtr interface{}) bool {
	mac := payloadPtr.(*layers.Ethernet)

	for _, field := range protocol.Fields {
		if field.NameId == parserRuler.EthSrcMACField {
			if !bytesEqual(mac.SrcMAC, field.ValuePtr.(net.HardwareAddr)) {
				return false
			}
		} else if field.NameId == parserRuler.EthDstMACField {
			if !bytesEqual(mac.DstMAC, field.ValuePtr.(net.HardwareAddr)) {
				return false
			}
		} else if field.NameId == parserRuler.EthLengthField {
			if mac.Length != field.ValuePtr.(uint16) {
				return false
			}
		} else if field.NameId == parserRuler.EthTypeField {
			if mac.EthernetType != field.ValuePtr.(layers.EthernetType) {
				return false
			}
		}
	}
	return true
}
