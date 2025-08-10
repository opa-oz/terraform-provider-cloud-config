package ccmodules

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

type CCModuleFlat struct {
	attributes map[string]schema.Attribute
}

func (cc *CCModuleFlat) Attributes() map[string]schema.Attribute {
	return cc.attributes
}

type CCModuleNested struct {
	block map[string]schema.Block
}

func (cc *CCModuleNested) Block() map[string]schema.Block {
	return cc.block
}
