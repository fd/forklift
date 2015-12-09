// http://sentry.readthedocs.org/en/latest/developer/client/index.html

package sentry

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"strings"
	"time"
)

var (
	DefaultPacket   Packet
	defaultLogger   = "root"
	defaultPlatform = fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
)

func init() {
	DefaultPacket.init()
}

type Packet struct {
	EventID    string                 `json:"event_id,omitempty"`
	Message    string                 `json:"message,omitempty"` // max 1000
	Timestamp  time.Time              `json:"timestamp,omitempty"`
	Level      LogLevel               `json:"level,omitempty"`
	Logger     string                 `json:"logger,omitempty"`
	Platform   string                 `json:"platform,omitempty"`
	Culprit    string                 `json:"culprit,omitempty"`
	Tags       Tags                   `json:"tags,omitempty"`
	ServerName string                 `json:"server_name,omitempty"`
	Modules    []map[string]string    `json:"modules,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`

	Stacktrace struct {
		Frames []*stack_frame_t `json:"frames"`
	} `json:"stacktrace,omitempty"`
}

/*

type stack_frame_t struct {
  Filename    string   `json:"filename"`
  Function    string   `json:"function"`
  Module      string   `json:"module"`
  Line        int      `json:"lineno"`
  AbsPath     string   `json:"abs_path"`
  PreContext  []string `json:"pre_context"`
  ContextLine string   `json:"context_line"`
  PostContext []string `json:"post_context"`
  InApp       bool     `json:"in_app"`
}

  for i, j := 0, len(stack)/2; i < j; i++ {
    k := len(stack) - i - 1
    stack[i], stack[k] = stack[k], stack[i]
  }
*/

type LogLevel uint8

type Tags [][2]string

const (
	FATAL   = LogLevel(1)
	ERROR   = LogLevel(2)
	WARNING = LogLevel(3)
	INFO    = LogLevel(4)
	DEBUG   = LogLevel(5)
)

func (level LogLevel) MarshalJSON() ([]byte, error) {
	switch level {
	case FATAL:
		return []byte(`"fatal"`), nil
	case ERROR:
		return []byte(`"error"`), nil
	case WARNING:
		return []byte(`"warning"`), nil
	case INFO:
		return []byte(`"info"`), nil
	case DEBUG:
		return []byte(`"debug"`), nil
	default:
		return []byte(`"error"`), nil
	}
}

func NewPacket() *Packet {
	p := DefaultPacket.copy()
	p.EventID = uuid4_str()
	p.Timestamp = time.Now().UTC()
	return p
}

func (packet *Packet) copy() *Packet {
	p := &Packet{}
	*p = *packet

	p.Tags = make([][2]string, len(p.Tags))
	for i, tag := range packet.Tags {
		p.Tags[i] = tag
	}

	return p
}

func (packet *Packet) init() {
	packet.Level = ERROR
	packet.Logger = defaultLogger
	packet.Platform = defaultPlatform
	packet.Tags.Add("os", runtime.GOOS)
	packet.Tags.Add("arch", runtime.GOARCH)
	packet.Tags.Add("go-version", runtime.Version())
	packet.Extra = make(map[string]interface{}, 1)
}

func (p *Packet) to_json() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Packet) Send() error {
	return call(p)
}

func (p *Packet) CaptureStack() {
	p.Stacktrace.Frames = stack()
}

func uuid4_str() string {
	var (
		u [16]uint8
	)

	rand.Read(u[:])

	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number from Section 4.1.3.
	u[6] = u[6]&0x0f | 0x40

	// Set the two most significant bits (bits 6 and 7) of the
	// clock_seq_hi_and_reserved to zero and one, respectively.
	u[8] = u[8]&0x3f | 0x80

	return hex.EncodeToString(u[:])
}

func (tags *Tags) Add(key, value string) {
	s := *tags
	s = append(s, [2]string{key, value})
	*tags = s
}
