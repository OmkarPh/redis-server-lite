package main

type Engine struct {
	// kv_store map[string]interface{}
	kv_store map[string]string
}

func (engine *Engine) get(key string) (string, bool) {
	if value, ok := engine.kv_store[key]; ok {
		return value, true
	}
	return "", false
}

func (engine *Engine) set(key string, value string) {
	engine.kv_store[key] = value
}

var KvEngine = Engine{
	kv_store: map[string]string{
		"dev": "Omkar Phansopkar",
	},
}
