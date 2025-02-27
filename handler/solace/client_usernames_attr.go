package solace

import (
	"context"
	"encoding/json"
	"strings"

	scalablev1alpha1 "github.com/benm-stm/solace-scalable-k8s-operator/api/v1alpha1"
	handler "github.com/benm-stm/solace-scalable-k8s-operator/handler"
)

// ClientUsername Attribute response struct
type ClientUsernameAttribute struct {
	AttributeName  string `json:"attributeName"`
	AttributeValue string `json:"attributeValue"`
	ClientUsername string `json:"clientUsername"`
	MsgVpnName     string `json:"msgVpnName"`
}

// ClientUsername Attributes array response struct
type ClientUsernameAttributes struct {
	Data []ClientUsernameAttribute `json:"data"`
}

var attrPath = "/config/msgVpns/" +
	"<MsgVpnName>" +
	"/clientUsernames/" +
	"<ClientUsername>" +
	"/attributes"

func NewClientUsernameAttrs() *ClientUsernameAttributes {
	return &ClientUsernameAttributes{}
}

// Returns the solace's clientUsername attributes in Json format
func (c *ClientUsernameAttributes) Add(
	s *scalablev1alpha1.SolaceScalable,
	n int,
	m *ClientUsername,
	p string,
	ctx context.Context,
) error {
	path := strings.Replace(attrPath, "<MsgVpnName>", m.MsgVpnName, 1)
	path = strings.Replace(path, "<ClientUsername>", m.ClientUsername, 1)
	url := handler.ConstructSempUrl(*s, n, path, nil)
	textBytes, _, _ := handler.CallSolaceSempApi(
		url,
		ctx,
		p,
	)
	resp := NewClientUsernameAttrs()
	if err := json.Unmarshal(textBytes, &resp); err != nil {
		return err
	}
	c.Data = append(
		c.Data,
		resp.Data...,
	)
	return nil
}

/*
func (c *clientUsernameAttributes) GetClientUsernameAttributes() []clientUsernameAttribute {
	return c.Data
}

func (c *clientUsernameAttribute) GetClientUsernameAttribute() *clientUsernameAttribute {
	return c
}

func (c *clientUsernameAttribute) GetMsgVpnName() string {
	return c.MsgVpnName
}

func (c *clientUsernameAttribute) GetClientUsername() string {
	return c.ClientUsername
}

func (c *clientUsernameAttribute) GetAttributeName() string {
	return c.AttributeName
}

func (c *clientUsernameAttribute) GetAttributeValue() string {
	return c.AttributeValue
}
*/
