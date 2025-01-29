package hclhelpers

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func GetFirstBlockOfType(blocks hcl.Blocks, blockType string) *hcl.Block {
	for _, block := range blocks {
		if block.Type == blockType {
			return block
		}
	}
	return nil
}

func FindChildBlocks(parentBlock *hcl.Block, blockType string) hcl.Blocks {
	var res hcl.Blocks
	childBlocks := parentBlock.Body.(*hclsyntax.Body).Blocks
	for _, b := range childBlocks {
		if b.Type == blockType {
			res = append(res, b.AsHCLBlock())
		}
	}
	return res
}
func FindFirstChildBlock(parentBlock *hcl.Block, blockType string) *hcl.Block {
	childBlocks := FindChildBlocks(parentBlock, blockType)
	if len(childBlocks) == 0 {
		return nil
	}
	return childBlocks[0]
}

// BlocksToMap convert an array of blocks to a map keyed by block laabel
// NOTE: this panics if any blocks do not have a label
func BlocksToMap(blocks hcl.Blocks) map[string]*hcl.Block {
	res := make(map[string]*hcl.Block, len(blocks))
	for _, b := range blocks {
		if len(b.Labels) == 0 {
			panic("all blocks passed to BlocksToMap must have a label")
		}
		res[b.Labels[0]] = b
	}
	return res
}

// BlockRange returns the range for the block
func BlockRange(block *hcl.Block) hcl.Range {
	if hclBody, ok := block.Body.(*hclsyntax.Body); ok {
		return hclBody.SrcRange
	}
	return block.DefRange
}

// BlockRangeWithLabels returns the range for the block including the labels
func BlockRangeWithLabels(block *hcl.Block) hcl.Range {
	if hclBody, ok := block.Body.(*hclsyntax.Body); ok {
		// build a range including the block labels
		var res hcl.Range
		res.Filename = hclBody.SrcRange.Filename
		res.Start = block.DefRange.Start
		res.End = hclBody.SrcRange.End
		return res
	}
	return block.DefRange
}

func BlockRangePointer(block *hcl.Block) *hcl.Range {
	r := BlockRange(block)
	return &r
}
func HclSyntaxBlockRangePointer(block *hclsyntax.Block) *hcl.Range {
	r := block.Range()
	return &r
}
