package selector

type Rounder interface {
	Get(key string) string
	ReSet(rules []string)
	Init(rules []string)
	CurrentNode() []string
}
