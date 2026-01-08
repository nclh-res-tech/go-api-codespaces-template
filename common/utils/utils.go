package utils

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractSOAPPayload pulls out the payload under the given root element from a SOAP envelope.
func ExtractSOAPPayload(body []byte, targetRoot string) []byte {
	if len(body) == 0 {
		return body
	}
	decoded := html.UnescapeString(string(body))
	if targetRoot != "" {
		tag := "<" + targetRoot
		if idx := strings.Index(decoded, tag); idx >= 0 {
			return []byte(decoded[idx:])
		}
	}
	if idx := strings.Index(decoded, "<"); idx >= 0 {
		return []byte(decoded[idx:])
	}
	return []byte(decoded)
}

// ErrString safely stringifies an error.
func ErrString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// XMLRequestToJSON reads the request body, unmarshals into T, applies setRoot, and responds with JSON.
func XMLRequestToJSON[T any](c *gin.Context, setRoot func(*T)) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty body"})
		return
	}

	var obj T
	if err := xml.Unmarshal(body, &obj); err == nil {
		if setRoot != nil {
			setRoot(&obj)
		}
		if generic, err := xmlToMap(body); err == nil {
			c.JSON(http.StatusOK, chooseRichest(obj, generic))
			return
		}
		c.JSON(http.StatusOK, obj)
		return
	}

	if generic, err := xmlToMap(body); err == nil {
		c.JSON(http.StatusOK, generic)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode XML"})
}

// XMLToJSON tries to decode the body as request or response types and returns the first successful JSON conversion.
func XMLToJSON[TReq any, TRes any](c *gin.Context, setReq func(*TReq), setRes func(*TRes)) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(body) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty body"})
		return
	}

	var lastErr error
	typeHint := strings.ToLower(strings.TrimSpace(c.Query("type")))
	if typeHint == "" || typeHint == "request" {
		var obj TReq
		if err := xml.Unmarshal(body, &obj); err == nil {
			if setReq != nil {
				setReq(&obj)
			}
			if generic, err := xmlToMap(body); err == nil {
				c.JSON(http.StatusOK, chooseRichest(obj, generic))
				return
			}
			c.JSON(http.StatusOK, obj)
			return
		}
		lastErr = err
		if typeHint == "request" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if typeHint == "response" || typeHint == "" {
		var obj TRes
		if err := xml.Unmarshal(body, &obj); err == nil {
			if setRes != nil {
				setRes(&obj)
			}
			if generic, err := xmlToMap(body); err == nil {
				c.JSON(http.StatusOK, chooseRichest(obj, generic))
				return
			}
			c.JSON(http.StatusOK, obj)
			return
		}
		lastErr = err
		if typeHint == "response" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}
	if generic, err := xmlToMap(body); err == nil {
		c.JSON(http.StatusOK, generic)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": "failed to decode XML", "detail": ErrString(lastErr)})
}

// chooseRichest returns the payload with more data (leaf count wins).
func chooseRichest(typed any, generic map[string]any) any {
	if generic == nil {
		return typed
	}
	typedMap := map[string]any{}
	if b, err := json.Marshal(typed); err == nil {
		_ = json.Unmarshal(b, &typedMap)
	}
	if leafCount(generic) >= leafCount(typedMap) {
		return generic
	}
	return typed
}

func leafCount(v any) int {
	switch val := v.(type) {
	case map[string]any:
		total := 0
		for _, vv := range val {
			total += leafCount(vv)
		}
		return total
	case []any:
		total := 0
		for _, vv := range val {
			total += leafCount(vv)
		}
		return total
	default:
		if val == nil {
			return 0
		}
		return 1
	}
}

// XMLBytesToMap exposes xmlToMap and unwraps the root element for easier merging.
func XMLBytesToMap(data []byte) (map[string]any, error) {
	m, err := xmlToMap(data)
	if err != nil {
		return nil, err
	}
	if len(m) == 1 {
		for _, v := range m {
			if inner, ok := v.(map[string]any); ok {
				return inner, nil
			}
		}
	}
	return m, nil
}

// xmlToMap converts XML bytes to a map preserving attributes, text, and repeated children.
func xmlToMap(data []byte) (map[string]any, error) {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	type node struct {
		name     string
		children map[string]any
		text     []string
	}
	var stack []node

	push := func(n node) {
		stack = append(stack, n)
	}
	pop := func() node {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return n
	}

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			n := node{name: toSnake(el.Name.Local), children: make(map[string]any)}
			for _, attr := range el.Attr {
				key := toSnake(attr.Name.Local)
				n.children[key] = strings.TrimSpace(attr.Value)
			}
			push(n)
		case xml.CharData:
			if len(stack) == 0 {
				continue
			}
			txt := strings.TrimSpace(string(el))
			if txt != "" {
				stack[len(stack)-1].text = append(stack[len(stack)-1].text, txt)
			}
		case xml.EndElement:
			if len(stack) == 0 {
				continue
			}
			n := pop()
			if len(n.text) > 0 {
				n.children["value"] = strings.Join(n.text, " ")
			}
			obj := n.children
			if len(obj) == 1 {
				if v, ok := obj["value"]; ok {
					obj = map[string]any{"value": v}
				}
			}
			if len(stack) == 0 {
				return map[string]any{toSnake(el.Name.Local): obj}, nil
			}
			parent := &stack[len(stack)-1]
			key := toSnake(el.Name.Local)
			if existing, ok := parent.children[key]; ok {
				switch e := existing.(type) {
				case []any:
					parent.children[key] = append(e, obj)
				default:
					parent.children[key] = []any{e, obj}
				}
			} else {
				parent.children[key] = obj
			}
		}
	}
	return nil, fmt.Errorf("no elements")
}

// toSnake converts CamelCase/PascalCase/dashed names to snake_case.
func toSnake(s string) string {
	if s == "" {
		return s
	}
	var out []rune
	for i, r := range s {
		if r == '-' || r == ' ' || r == '.' {
			out = append(out, '_')
			continue
		}
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				out = append(out, '_')
			}
			out = append(out, r+('a'-'A'))
			continue
		}
		out = append(out, r)
	}
	res := string(out)
	for strings.Contains(res, "__") {
		res = strings.ReplaceAll(res, "__", "_")
	}
	return strings.Trim(res, "_")
}
