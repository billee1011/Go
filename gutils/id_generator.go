package gutils

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

//epoch 是指定为1970年一月一日凌晨零点零分零秒，格林威治时间
//0 - 0000000000 0000000000 0000000000 0000000000 0 - 0000000000 - 0000000000 00
//第一位为未使用，接下来的41位为毫秒级时间(41位的长度可以使用69年)，
//然后10位的节点id，支持0-1023
//最后12位是毫秒内的计数（12位的计数顺序号支持每个节点每毫秒产生4096个ID序号）
//生成的ID整体上按照时间自增排序，并且整个分布式系统内不会产生ID碰撞（由nodeId分区）。经测试每秒能够产生26万个ID。

var (
	epoch     int64 = 1288834974657
	nodeBits  uint8 = 10
	stepBits  uint8 = 12
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	nodeMask        = nodeMax << stepBits
	stepMask  int64 = -1 ^ (-1 << stepBits)
	timeShift       = nodeBits + stepBits
	nodeShift       = stepBits
)

const encodeBase32Map = "ybndrfg8ejkmcpqxot1uwisza345h769"

var decodeBase32Map [256]byte

const encodeBase58Map = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"

var decodeBase58Map [256]byte

// JSONSyntaxError json语法错误
type JSONSyntaxError struct{ original []byte }

func (j JSONSyntaxError) Error() string {
	return fmt.Sprintf("invalid snowflake ID %q", string(j.original))
}

func init() {

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase58Map); i++ {
		decodeBase58Map[encodeBase58Map[i]] = byte(i)
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[i] = 0xFF
	}

	for i := 0; i < len(encodeBase32Map); i++ {
		decodeBase32Map[encodeBase32Map[i]] = byte(i)
	}
}

// ErrInvalidBase58 不合法的base58
var ErrInvalidBase58 = errors.New("invalid base58")

// ErrInvalidBase32 不合法的base32
var ErrInvalidBase32 = errors.New("invalid base32")

// Node node结构
type Node struct {
	mu   sync.Mutex
	time int64
	node int64
	step int64
}

// ID id类型转化
type ID int64

// NewNode 根据节点id初始化id生成器
func NewNode(node int64) (*Node, error) {
	nodeMax = -1 ^ (-1 << nodeBits)
	nodeMask = nodeMax << stepBits
	stepMask = -1 ^ (-1 << stepBits)
	timeShift = nodeBits + stepBits
	nodeShift = stepBits

	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}

	return &Node{
		time: 0,
		node: node,
		step: 0,
	}, nil
}

// Generate id生成方法
func (n *Node) Generate() ID {
	n.mu.Lock()
	now := time.Now().UnixNano() / 1000000

	if n.time == now {
		n.step = (n.step + 1) & stepMask

		if n.step == 0 {
			for now <= n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now
	r := ID((now-epoch)<<timeShift |
		(n.node << nodeShift) |
		(n.step),
	)

	n.mu.Unlock()
	return r
}

// Int64 转换成int64
func (f ID) Int64() int64 {
	return int64(f)
}

// String 转换成string
func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

//Base2 转换成 base2
func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}

// Base36 转换成base36
func (f ID) Base36() string {
	return strconv.FormatInt(int64(f), 36)
}

// Base32 转换成Base32
func (f ID) Base32() string {

	if f < 32 {
		return string(encodeBase32Map[f])
	}

	b := make([]byte, 0, 12)
	for f >= 32 {
		b = append(b, encodeBase32Map[f%32])
		f /= 32
	}
	b = append(b, encodeBase32Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

//ParseBase32 将id解析成base32位编码
func ParseBase32(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase32Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase32
		}
		id = id*32 + int64(decodeBase32Map[b[i]])
	}

	return ID(id), nil
}

//Base58 转换成base58
func (f ID) Base58() string {

	if f < 58 {
		return string(encodeBase58Map[f])
	}

	b := make([]byte, 0, 11)
	for f >= 58 {
		b = append(b, encodeBase58Map[f%58])
		f /= 58
	}
	b = append(b, encodeBase58Map[f])

	for x, y := 0, len(b)-1; x < y; x, y = x+1, y-1 {
		b[x], b[y] = b[y], b[x]
	}

	return string(b)
}

//ParseBase58 将id解析成base58位编码
func ParseBase58(b []byte) (ID, error) {

	var id int64

	for i := range b {
		if decodeBase58Map[b[i]] == 0xFF {
			return -1, ErrInvalidBase58
		}
		id = id*58 + int64(decodeBase58Map[b[i]])
	}

	return ID(id), nil
}

// Base64 转成base64
func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

//Bytes 转成byte
func (f ID) Bytes() []byte {
	return []byte(f.String())
}

//IntBytes 转成intbyte
func (f ID) IntBytes() [8]byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(f))
	return b
}

//Time 获取时间位
func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + epoch
}

//Node 获取节点位
func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

//Step 获取步长位
func (f ID) Step() int64 {
	return int64(f) & stepMask
}

//MarshalJSON json化
func (f ID) MarshalJSON() ([]byte, error) {
	buff := make([]byte, 0, 22)
	buff = append(buff, '"')
	buff = strconv.AppendInt(buff, int64(f), 10)
	buff = append(buff, '"')
	return buff, nil
}

//UnmarshalJSON 反json化
func (f *ID) UnmarshalJSON(b []byte) error {
	if len(b) < 3 || b[0] != '"' || b[len(b)-1] != '"' {
		return JSONSyntaxError{b}
	}

	i, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}

	*f = ID(i)
	return nil
}
