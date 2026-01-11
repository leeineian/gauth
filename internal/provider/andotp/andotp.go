package andotp

import (
	"encoding/json"
	"os"

	ga "github.com/grijul/go-andotp/andotp"
	"github.com/leeineian/gauth/internal/model"
)

type andotpNode struct {
	Secret         string        `json:"secret"`
	Issuer         string        `json:"issuer"`
	Label          string        `json:"label"`
	Digits         int           `json:"digits"`
	Type           string        `json:"type"`
	Algorithm      string        `json:"alogrithm"`
	Thumbnail      string        `json:"thumbnail"`
	Last_used      int64         `json:"last_used"`
	Used_frequency int           `json:"used_frequency"`
	Period         int           `json:"period"`
	Tags           []interface{} `json:"tags"`
}

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Import(filePath string, password string) ([]model.Account, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var nodes []andotpNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		if password != "" {
			decrypted, err := ga.Decrypt(data, password)
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(decrypted, &nodes); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	accounts := make([]model.Account, 0, len(nodes))
	for _, n := range nodes {
		accounts = append(accounts, model.Account{
			Secret:    n.Secret,
			Issuer:    n.Issuer,
			Label:     n.Label,
			Digits:    n.Digits,
			Type:      model.OTPType(n.Type),
			Algorithm: n.Algorithm,
			Period:    int64(n.Period),
			Misc: map[string]interface{}{
				"thumbnail":      n.Thumbnail,
				"last_used":      n.Last_used,
				"used_frequency": n.Used_frequency,
				"tags":           n.Tags,
			},
		})
	}

	return accounts, nil
}

func (p *Provider) Export(accounts []model.Account, password string) ([]byte, error) {
	nodes := make([]andotpNode, 0, len(accounts))
	for _, a := range accounts {
		node := andotpNode{
			Secret:    a.Secret,
			Issuer:    a.Issuer,
			Label:     a.Label,
			Digits:    a.Digits,
			Type:      string(a.Type),
			Algorithm: a.Algorithm,
			Period:    int(a.Period),
		}
		if a.Misc != nil {
			if v, ok := a.Misc["thumbnail"].(string); ok {
				node.Thumbnail = v
			}
			if v, ok := a.Misc["last_used"].(int64); ok {
				node.Last_used = v
			}
			if v, ok := a.Misc["used_frequency"].(int); ok {
				node.Used_frequency = v
			}
			if v, ok := a.Misc["tags"].([]interface{}); ok {
				node.Tags = v
			}
		}
		nodes = append(nodes, node)
	}

	data, err := json.Marshal(nodes)
	if err != nil {
		return nil, err
	}

	if password != "" {
		return ga.Encrypt(data, password)
	}

	return data, nil
}
