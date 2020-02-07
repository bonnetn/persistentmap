package internal

type NodeMetadata struct {
	level int32
}

func (n NodeMetadata) Level() int32 {
	return n.level
}

func (n NodeMetadata) Child() NodeMetadata {
	return NodeMetadata{level: n.level + 1}
}
