package paxoi

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"io"
	"log"
	"sync"

	"github.com/vonaka/shreplic/state"
	"github.com/vonaka/shreplic/tools/fastrpc"
)

//////////////////////////////////////////////////////////////
//                                                          //
//  gobin-codegen doesn't support declarations of the form  //
//                                                          //
//      type A B                                            //
//                                                          //
//  that's why we use `[]CommandId` instead of `Dep`        //
//                                                          //
//////////////////////////////////////////////////////////////

type MFastAck struct {
	Replica  int32
	Ballot   int32
	CmdId    CommandId
	Dep      []CommandId
	Checksum []SHash
}

type MSlowAck struct {
	Replica  int32
	Ballot   int32
	CmdId    CommandId
	Dep      []CommandId
	Checksum []SHash
}

type MLightSlowAck struct {
	Replica int32
	Ballot  int32
	CmdId   CommandId
}

type MAcks struct {
	FastAcks      []MFastAck
	LightSlowAcks []MLightSlowAck
}

type Ack struct {
	CmdId    CommandId
	Dep      []CommandId
	Checksum []SHash
}

type MOptAcks struct {
	Replica int32
	Ballot  int32
	Acks    []Ack
}

type MReply struct {
	Replica  int32
	Ballot   int32
	CmdId    CommandId
	Checksum []SHash
	Rep      []byte
}

type MReadReply struct {
	Replica int32
	Ballot  int32
	CmdId   CommandId
	Rep     []byte
}

type MNewLeader struct {
	Replica int32
	Ballot  int32
}

type MNewLeaderAck struct {
	Replica int32
	Ballot  int32
	Cballot int32
}

type MShareState struct {
	Replica int32
	Ballot  int32
}

type MSync struct {
	Replica int32
	Ballot  int32
	Phases  map[CommandId]int
	Cmds    map[CommandId]state.Command
	Deps    map[CommandId]Dep
}

type MLightSync struct {
	Replica int32
	Ballot  int32
}

type MCollect struct {
	Replica int32
	Ballot  int32
	Ids     []CommandId
}

type MPing struct {
	Replica int32
	Ballot  int32
}

type MPingRep struct {
	Replica int32
	Ballot  int32
}

var (
	useFastAckPool = false
	fastAckPool    = sync.Pool{
		New: func() interface{} {
			return &MFastAck{}
		},
	}
)

func newFastAck() *MFastAck {
	if useFastAckPool {
		return fastAckPool.Get().(*MFastAck)
	}
	return &MFastAck{}
}

func releaseFastAck(f *MFastAck) {
	if useFastAckPool {
		fastAckPool.Put(f)
	}
}

func copyFastAck(fa *MFastAck) *MFastAck {
	fa2 := newFastAck()
	fa2.Replica = fa.Replica
	fa2.Ballot = fa.Ballot
	fa2.CmdId = fa.CmdId
	fa2.Dep = fa.Dep
	fa2.Checksum = fa.Checksum
	return fa2
}

func (m *MFastAck) New() fastrpc.Serializable {
	return newFastAck()
}

func (m *MSlowAck) New() fastrpc.Serializable {
	return (*MSlowAck)(newFastAck())
}

func (m *MLightSlowAck) New() fastrpc.Serializable {
	return new(MLightSlowAck)
}

func (m *MAcks) New() fastrpc.Serializable {
	return new(MAcks)
}

func (m *MOptAcks) New() fastrpc.Serializable {
	return new(MOptAcks)
}

func (m *MReply) New() fastrpc.Serializable {
	return new(MReply)
}

func (m *MReadReply) New() fastrpc.Serializable {
	return new(MReadReply)
}

func (m *MNewLeader) New() fastrpc.Serializable {
	return new(MNewLeader)
}

func (m *MNewLeaderAck) New() fastrpc.Serializable {
	return new(MNewLeaderAck)
}

func (m *MShareState) New() fastrpc.Serializable {
	return new(MShareState)
}

func (m *MSync) New() fastrpc.Serializable {
	return new(MSync)
}

func (m *MSync) Marshal(w io.Writer) {
	e := gob.NewEncoder(w)

	err := e.Encode(m)
	if err != nil {
		log.Fatal("Don't know what to do")
	}
}

func (m *MSync) Unmarshal(r io.Reader) error {
	e := gob.NewDecoder(r)

	return e.Decode(m)
}

func (m *MLightSync) New() fastrpc.Serializable {
	return new(MLightSync)
}

func (m *MCollect) New() fastrpc.Serializable {
	return new(MCollect)
}

func (m *MPing) New() fastrpc.Serializable {
	return new(MPing)
}

func (m *MPingRep) New() fastrpc.Serializable {
	return new(MPingRep)
}

///////////////////////////////////////////////////////////////////////////////
//                                                                           //
//  Generated with gobin-codegen [https://code.google.com/p/gobin-codegen/]  //
//                                                                           //
///////////////////////////////////////////////////////////////////////////////

type byteReader interface {
	io.Reader
	ReadByte() (c byte, err error)
}

func (t *CommandId) BinarySize() (nbytes int, sizeKnown bool) {
	return 8, true
}

type CommandIdCache struct {
	mu	sync.Mutex
	cache	[]*CommandId
}

func NewCommandIdCache() *CommandIdCache {
	c := &CommandIdCache{}
	c.cache = make([]*CommandId, 0)
	return c
}

func (p *CommandIdCache) Get() *CommandId {
	var t *CommandId
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &CommandId{}
	}
	return t
}
func (p *CommandIdCache) Put(t *CommandId) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *CommandId) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.ClientId
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.SeqNum
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *CommandId) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.SeqNum = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	return nil
}

func (t *SHash) BinarySize() (nbytes int, sizeKnown bool) {
	return 32, true
}

type SHashCache struct {
	mu	sync.Mutex
	cache	[]*SHash
}

func NewSHashCache() *SHashCache {
	c := &SHashCache{}
	c.cache = make([]*SHash, 0)
	return c
}

func (p *SHashCache) Get() *SHash {
	var t *SHash
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &SHash{}
	}
	return t
}
func (p *SHashCache) Put(t *SHash) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *SHash) Marshal(wire io.Writer) {
	var b [32]byte
	var bs []byte
	bs = b[:32]
	bs[0] = byte(t.H[0])
	bs[1] = byte(t.H[1])
	bs[2] = byte(t.H[2])
	bs[3] = byte(t.H[3])
	bs[4] = byte(t.H[4])
	bs[5] = byte(t.H[5])
	bs[6] = byte(t.H[6])
	bs[7] = byte(t.H[7])
	bs[8] = byte(t.H[8])
	bs[9] = byte(t.H[9])
	bs[10] = byte(t.H[10])
	bs[11] = byte(t.H[11])
	bs[12] = byte(t.H[12])
	bs[13] = byte(t.H[13])
	bs[14] = byte(t.H[14])
	bs[15] = byte(t.H[15])
	bs[16] = byte(t.H[16])
	bs[17] = byte(t.H[17])
	bs[18] = byte(t.H[18])
	bs[19] = byte(t.H[19])
	bs[20] = byte(t.H[20])
	bs[21] = byte(t.H[21])
	bs[22] = byte(t.H[22])
	bs[23] = byte(t.H[23])
	bs[24] = byte(t.H[24])
	bs[25] = byte(t.H[25])
	bs[26] = byte(t.H[26])
	bs[27] = byte(t.H[27])
	bs[28] = byte(t.H[28])
	bs[29] = byte(t.H[29])
	bs[30] = byte(t.H[30])
	bs[31] = byte(t.H[31])
	wire.Write(bs)
}

func (t *SHash) Unmarshal(wire io.Reader) error {
	var b [32]byte
	var bs []byte
	bs = b[:32]
	if _, err := io.ReadAtLeast(wire, bs, 32); err != nil {
		return err
	}
	t.H[0] = byte(bs[0])
	t.H[1] = byte(bs[1])
	t.H[2] = byte(bs[2])
	t.H[3] = byte(bs[3])
	t.H[4] = byte(bs[4])
	t.H[5] = byte(bs[5])
	t.H[6] = byte(bs[6])
	t.H[7] = byte(bs[7])
	t.H[8] = byte(bs[8])
	t.H[9] = byte(bs[9])
	t.H[10] = byte(bs[10])
	t.H[11] = byte(bs[11])
	t.H[12] = byte(bs[12])
	t.H[13] = byte(bs[13])
	t.H[14] = byte(bs[14])
	t.H[15] = byte(bs[15])
	t.H[16] = byte(bs[16])
	t.H[17] = byte(bs[17])
	t.H[18] = byte(bs[18])
	t.H[19] = byte(bs[19])
	t.H[20] = byte(bs[20])
	t.H[21] = byte(bs[21])
	t.H[22] = byte(bs[22])
	t.H[23] = byte(bs[23])
	t.H[24] = byte(bs[24])
	t.H[25] = byte(bs[25])
	t.H[26] = byte(bs[26])
	t.H[27] = byte(bs[27])
	t.H[28] = byte(bs[28])
	t.H[29] = byte(bs[29])
	t.H[30] = byte(bs[30])
	t.H[31] = byte(bs[31])
	return nil
}

func (t *MAcks) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MAcksCache struct {
	mu	sync.Mutex
	cache	[]*MAcks
}

func NewMAcksCache() *MAcksCache {
	c := &MAcksCache{}
	c.cache = make([]*MAcks, 0)
	return c
}

func (p *MAcksCache) Get() *MAcks {
	var t *MAcks
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MAcks{}
	}
	return t
}
func (p *MAcksCache) Put(t *MAcks) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MAcks) Marshal(wire io.Writer) {
	var b [10]byte
	var bs []byte
	bs = b[:]
	alen1 := int64(len(t.FastAcks))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		t.FastAcks[i].Marshal(wire)
	}
	bs = b[:]
	alen2 := int64(len(t.LightSlowAcks))
	if wlen := binary.PutVarint(bs, alen2); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen2; i++ {
		bs = b[:4]
		tmp32 := t.LightSlowAcks[i].Replica
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.LightSlowAcks[i].Ballot
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.LightSlowAcks[i].CmdId.ClientId
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.LightSlowAcks[i].CmdId.SeqNum
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
	}
}

func (t *MAcks) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [10]byte
	var bs []byte
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.FastAcks = make([]MFastAck, alen1)
	for i := int64(0); i < alen1; i++ {
		t.FastAcks[i].Unmarshal(wire)
	}
	alen2, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.LightSlowAcks = make([]MLightSlowAck, alen2)
	for i := int64(0); i < alen2; i++ {
		bs = b[:4]
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.LightSlowAcks[i].Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.LightSlowAcks[i].Ballot = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.LightSlowAcks[i].CmdId.ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.LightSlowAcks[i].CmdId.SeqNum = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	}
	return nil
}

func (t *MNewLeader) BinarySize() (nbytes int, sizeKnown bool) {
	return 8, true
}

type MNewLeaderCache struct {
	mu	sync.Mutex
	cache	[]*MNewLeader
}

func NewMNewLeaderCache() *MNewLeaderCache {
	c := &MNewLeaderCache{}
	c.cache = make([]*MNewLeader, 0)
	return c
}

func (p *MNewLeaderCache) Get() *MNewLeader {
	var t *MNewLeader
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MNewLeader{}
	}
	return t
}
func (p *MNewLeaderCache) Put(t *MNewLeader) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MNewLeader) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MNewLeader) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	return nil
}

func (t *Ack) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type AckCache struct {
	mu	sync.Mutex
	cache	[]*Ack
}

func NewAckCache() *AckCache {
	c := &AckCache{}
	c.cache = make([]*Ack, 0)
	return c
}

func (p *AckCache) Get() *Ack {
	var t *Ack
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &Ack{}
	}
	return t
}
func (p *AckCache) Put(t *Ack) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *Ack) Marshal(wire io.Writer) {
	var b [10]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.CmdId.ClientId
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.SeqNum
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Dep))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		tmp32 = t.Dep[i].ClientId
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.Dep[i].SeqNum
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
	}
	bs = b[:]
	alen2 := int64(len(t.Checksum))
	if wlen := binary.PutVarint(bs, alen2); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen2; i++ {
		bs = b[:1]
		bs[0] = byte(t.Checksum[i].H[0])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[1])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[2])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[3])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[4])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[5])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[6])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[7])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[8])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[9])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[10])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[11])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[12])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[13])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[14])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[15])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[16])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[17])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[18])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[19])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[20])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[21])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[22])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[23])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[24])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[25])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[26])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[27])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[28])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[29])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[30])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[31])
		wire.Write(bs)
	}
}

func (t *Ack) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [10]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.CmdId.ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.CmdId.SeqNum = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Dep = make([]CommandId, alen1)
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Dep[i].ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Dep[i].SeqNum = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	}
	alen2, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Checksum = make([]SHash, alen2)
	for i := int64(0); i < alen2; i++ {
		bs = b[:1]
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[0] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[1] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[2] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[3] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[4] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[5] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[6] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[7] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[8] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[9] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[10] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[11] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[12] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[13] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[14] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[15] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[16] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[17] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[18] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[19] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[20] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[21] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[22] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[23] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[24] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[25] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[26] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[27] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[28] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[29] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[30] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[31] = byte(bs[0])
	}
	return nil
}

func (t *MOptAcks) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MOptAcksCache struct {
	mu	sync.Mutex
	cache	[]*MOptAcks
}

func NewMOptAcksCache() *MOptAcksCache {
	c := &MOptAcksCache{}
	c.cache = make([]*MOptAcks, 0)
	return c
}

func (p *MOptAcksCache) Get() *MOptAcks {
	var t *MOptAcks
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MOptAcks{}
	}
	return t
}
func (p *MOptAcksCache) Put(t *MOptAcks) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MOptAcks) Marshal(wire io.Writer) {
	var b [10]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Acks))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		t.Acks[i].Marshal(wire)
	}
}

func (t *MOptAcks) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [10]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Acks = make([]Ack, alen1)
	for i := int64(0); i < alen1; i++ {
		t.Acks[i].Unmarshal(wire)
	}
	return nil
}

func (t *MReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MReplyCache struct {
	mu	sync.Mutex
	cache	[]*MReply
}

func NewMReplyCache() *MReplyCache {
	c := &MReplyCache{}
	c.cache = make([]*MReply, 0)
	return c
}

func (p *MReplyCache) Get() *MReply {
	var t *MReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MReply{}
	}
	return t
}
func (p *MReplyCache) Put(t *MReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MReply) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.ClientId
	bs[8] = byte(tmp32)
	bs[9] = byte(tmp32 >> 8)
	bs[10] = byte(tmp32 >> 16)
	bs[11] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.SeqNum
	bs[12] = byte(tmp32)
	bs[13] = byte(tmp32 >> 8)
	bs[14] = byte(tmp32 >> 16)
	bs[15] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Checksum))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		bs = b[:1]
		bs[0] = byte(t.Checksum[i].H[0])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[1])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[2])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[3])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[4])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[5])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[6])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[7])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[8])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[9])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[10])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[11])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[12])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[13])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[14])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[15])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[16])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[17])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[18])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[19])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[20])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[21])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[22])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[23])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[24])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[25])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[26])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[27])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[28])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[29])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[30])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[31])
		wire.Write(bs)
	}
	bs = b[:]
	alen3 := int64(len(t.Rep))
	if wlen := binary.PutVarint(bs, alen3); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen3; i++ {
		bs = b[:1]
		bs[0] = byte(t.Rep[i])
		wire.Write(bs)
	}
}

func (t *MReply) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	t.CmdId.ClientId = int32((uint32(bs[8]) | (uint32(bs[9]) << 8) | (uint32(bs[10]) << 16) | (uint32(bs[11]) << 24)))
	t.CmdId.SeqNum = int32((uint32(bs[12]) | (uint32(bs[13]) << 8) | (uint32(bs[14]) << 16) | (uint32(bs[15]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Checksum = make([]SHash, alen1)
	for i := int64(0); i < alen1; i++ {
		bs = b[:1]
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[0] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[1] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[2] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[3] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[4] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[5] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[6] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[7] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[8] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[9] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[10] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[11] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[12] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[13] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[14] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[15] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[16] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[17] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[18] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[19] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[20] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[21] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[22] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[23] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[24] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[25] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[26] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[27] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[28] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[29] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[30] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[31] = byte(bs[0])
	}
	alen3, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Rep = make([]byte, alen3)
	for i := int64(0); i < alen3; i++ {
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Rep[i] = byte(bs[0])
	}
	return nil
}

func (t *MCollect) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MCollectCache struct {
	mu	sync.Mutex
	cache	[]*MCollect
}

func NewMCollectCache() *MCollectCache {
	c := &MCollectCache{}
	c.cache = make([]*MCollect, 0)
	return c
}

func (p *MCollectCache) Get() *MCollect {
	var t *MCollect
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MCollect{}
	}
	return t
}
func (p *MCollectCache) Put(t *MCollect) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MCollect) Marshal(wire io.Writer) {
	var b [10]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Ids))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		tmp32 = t.Ids[i].ClientId
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.Ids[i].SeqNum
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
	}
}

func (t *MCollect) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [10]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Ids = make([]CommandId, alen1)
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Ids[i].ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Ids[i].SeqNum = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	}
	return nil
}

func (t *MNewLeaderAck) BinarySize() (nbytes int, sizeKnown bool) {
	return 12, true
}

type MNewLeaderAckCache struct {
	mu	sync.Mutex
	cache	[]*MNewLeaderAck
}

func NewMNewLeaderAckCache() *MNewLeaderAckCache {
	c := &MNewLeaderAckCache{}
	c.cache = make([]*MNewLeaderAck, 0)
	return c
}

func (p *MNewLeaderAckCache) Get() *MNewLeaderAck {
	var t *MNewLeaderAck
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MNewLeaderAck{}
	}
	return t
}
func (p *MNewLeaderAckCache) Put(t *MNewLeaderAck) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MNewLeaderAck) Marshal(wire io.Writer) {
	var b [12]byte
	var bs []byte
	bs = b[:12]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	tmp32 = t.Cballot
	bs[8] = byte(tmp32)
	bs[9] = byte(tmp32 >> 8)
	bs[10] = byte(tmp32 >> 16)
	bs[11] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MNewLeaderAck) Unmarshal(wire io.Reader) error {
	var b [12]byte
	var bs []byte
	bs = b[:12]
	if _, err := io.ReadAtLeast(wire, bs, 12); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	t.Cballot = int32((uint32(bs[8]) | (uint32(bs[9]) << 8) | (uint32(bs[10]) << 16) | (uint32(bs[11]) << 24)))
	return nil
}

func (t *MShareState) BinarySize() (nbytes int, sizeKnown bool) {
	return 8, true
}

type MShareStateCache struct {
	mu	sync.Mutex
	cache	[]*MShareState
}

func NewMShareStateCache() *MShareStateCache {
	c := &MShareStateCache{}
	c.cache = make([]*MShareState, 0)
	return c
}

func (p *MShareStateCache) Get() *MShareState {
	var t *MShareState
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MShareState{}
	}
	return t
}
func (p *MShareStateCache) Put(t *MShareState) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MShareState) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MShareState) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	return nil
}

func (t *MLightSync) BinarySize() (nbytes int, sizeKnown bool) {
	return 8, true
}

type MLightSyncCache struct {
	mu	sync.Mutex
	cache	[]*MLightSync
}

func NewMLightSyncCache() *MLightSyncCache {
	c := &MLightSyncCache{}
	c.cache = make([]*MLightSync, 0)
	return c
}

func (p *MLightSyncCache) Get() *MLightSync {
	var t *MLightSync
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MLightSync{}
	}
	return t
}
func (p *MLightSyncCache) Put(t *MLightSync) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MLightSync) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MLightSync) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	return nil
}

func (t *MFastAck) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MFastAckCache struct {
	mu	sync.Mutex
	cache	[]*MFastAck
}

func NewMFastAckCache() *MFastAckCache {
	c := &MFastAckCache{}
	c.cache = make([]*MFastAck, 0)
	return c
}

func (p *MFastAckCache) Get() *MFastAck {
	var t *MFastAck
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MFastAck{}
	}
	return t
}
func (p *MFastAckCache) Put(t *MFastAck) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MFastAck) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.ClientId
	bs[8] = byte(tmp32)
	bs[9] = byte(tmp32 >> 8)
	bs[10] = byte(tmp32 >> 16)
	bs[11] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.SeqNum
	bs[12] = byte(tmp32)
	bs[13] = byte(tmp32 >> 8)
	bs[14] = byte(tmp32 >> 16)
	bs[15] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Dep))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		tmp32 = t.Dep[i].ClientId
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.Dep[i].SeqNum
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
	}
	bs = b[:]
	alen2 := int64(len(t.Checksum))
	if wlen := binary.PutVarint(bs, alen2); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen2; i++ {
		bs = b[:1]
		bs[0] = byte(t.Checksum[i].H[0])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[1])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[2])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[3])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[4])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[5])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[6])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[7])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[8])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[9])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[10])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[11])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[12])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[13])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[14])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[15])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[16])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[17])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[18])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[19])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[20])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[21])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[22])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[23])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[24])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[25])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[26])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[27])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[28])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[29])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[30])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[31])
		wire.Write(bs)
	}
}

func (t *MFastAck) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	t.CmdId.ClientId = int32((uint32(bs[8]) | (uint32(bs[9]) << 8) | (uint32(bs[10]) << 16) | (uint32(bs[11]) << 24)))
	t.CmdId.SeqNum = int32((uint32(bs[12]) | (uint32(bs[13]) << 8) | (uint32(bs[14]) << 16) | (uint32(bs[15]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Dep = make([]CommandId, alen1)
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Dep[i].ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Dep[i].SeqNum = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	}
	alen2, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Checksum = make([]SHash, alen2)
	for i := int64(0); i < alen2; i++ {
		bs = b[:1]
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[0] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[1] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[2] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[3] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[4] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[5] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[6] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[7] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[8] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[9] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[10] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[11] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[12] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[13] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[14] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[15] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[16] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[17] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[18] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[19] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[20] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[21] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[22] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[23] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[24] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[25] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[26] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[27] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[28] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[29] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[30] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[31] = byte(bs[0])
	}
	return nil
}

func (t *MSlowAck) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MSlowAckCache struct {
	mu	sync.Mutex
	cache	[]*MSlowAck
}

func NewMSlowAckCache() *MSlowAckCache {
	c := &MSlowAckCache{}
	c.cache = make([]*MSlowAck, 0)
	return c
}

func (p *MSlowAckCache) Get() *MSlowAck {
	var t *MSlowAck
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MSlowAck{}
	}
	return t
}
func (p *MSlowAckCache) Put(t *MSlowAck) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MSlowAck) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.ClientId
	bs[8] = byte(tmp32)
	bs[9] = byte(tmp32 >> 8)
	bs[10] = byte(tmp32 >> 16)
	bs[11] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.SeqNum
	bs[12] = byte(tmp32)
	bs[13] = byte(tmp32 >> 8)
	bs[14] = byte(tmp32 >> 16)
	bs[15] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Dep))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		tmp32 = t.Dep[i].ClientId
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
		tmp32 = t.Dep[i].SeqNum
		bs[0] = byte(tmp32)
		bs[1] = byte(tmp32 >> 8)
		bs[2] = byte(tmp32 >> 16)
		bs[3] = byte(tmp32 >> 24)
		wire.Write(bs)
	}
	bs = b[:]
	alen2 := int64(len(t.Checksum))
	if wlen := binary.PutVarint(bs, alen2); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen2; i++ {
		bs = b[:1]
		bs[0] = byte(t.Checksum[i].H[0])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[1])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[2])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[3])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[4])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[5])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[6])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[7])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[8])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[9])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[10])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[11])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[12])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[13])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[14])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[15])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[16])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[17])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[18])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[19])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[20])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[21])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[22])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[23])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[24])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[25])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[26])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[27])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[28])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[29])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[30])
		wire.Write(bs)
		bs[0] = byte(t.Checksum[i].H[31])
		wire.Write(bs)
	}
}

func (t *MSlowAck) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	t.CmdId.ClientId = int32((uint32(bs[8]) | (uint32(bs[9]) << 8) | (uint32(bs[10]) << 16) | (uint32(bs[11]) << 24)))
	t.CmdId.SeqNum = int32((uint32(bs[12]) | (uint32(bs[13]) << 8) | (uint32(bs[14]) << 16) | (uint32(bs[15]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Dep = make([]CommandId, alen1)
	for i := int64(0); i < alen1; i++ {
		bs = b[:4]
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Dep[i].ClientId = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
		if _, err := io.ReadAtLeast(wire, bs, 4); err != nil {
			return err
		}
		t.Dep[i].SeqNum = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	}
	alen2, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Checksum = make([]SHash, alen2)
	for i := int64(0); i < alen2; i++ {
		bs = b[:1]
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[0] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[1] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[2] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[3] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[4] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[5] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[6] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[7] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[8] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[9] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[10] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[11] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[12] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[13] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[14] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[15] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[16] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[17] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[18] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[19] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[20] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[21] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[22] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[23] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[24] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[25] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[26] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[27] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[28] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[29] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[30] = byte(bs[0])
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Checksum[i].H[31] = byte(bs[0])
	}
	return nil
}

func (t *MLightSlowAck) BinarySize() (nbytes int, sizeKnown bool) {
	return 16, true
}

type MLightSlowAckCache struct {
	mu	sync.Mutex
	cache	[]*MLightSlowAck
}

func NewMLightSlowAckCache() *MLightSlowAckCache {
	c := &MLightSlowAckCache{}
	c.cache = make([]*MLightSlowAck, 0)
	return c
}

func (p *MLightSlowAckCache) Get() *MLightSlowAck {
	var t *MLightSlowAck
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MLightSlowAck{}
	}
	return t
}
func (p *MLightSlowAckCache) Put(t *MLightSlowAck) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MLightSlowAck) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.ClientId
	bs[8] = byte(tmp32)
	bs[9] = byte(tmp32 >> 8)
	bs[10] = byte(tmp32 >> 16)
	bs[11] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.SeqNum
	bs[12] = byte(tmp32)
	bs[13] = byte(tmp32 >> 8)
	bs[14] = byte(tmp32 >> 16)
	bs[15] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MLightSlowAck) Unmarshal(wire io.Reader) error {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	t.CmdId.ClientId = int32((uint32(bs[8]) | (uint32(bs[9]) << 8) | (uint32(bs[10]) << 16) | (uint32(bs[11]) << 24)))
	t.CmdId.SeqNum = int32((uint32(bs[12]) | (uint32(bs[13]) << 8) | (uint32(bs[14]) << 16) | (uint32(bs[15]) << 24)))
	return nil
}

func (t *MReadReply) BinarySize() (nbytes int, sizeKnown bool) {
	return 0, false
}

type MReadReplyCache struct {
	mu	sync.Mutex
	cache	[]*MReadReply
}

func NewMReadReplyCache() *MReadReplyCache {
	c := &MReadReplyCache{}
	c.cache = make([]*MReadReply, 0)
	return c
}

func (p *MReadReplyCache) Get() *MReadReply {
	var t *MReadReply
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MReadReply{}
	}
	return t
}
func (p *MReadReplyCache) Put(t *MReadReply) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MReadReply) Marshal(wire io.Writer) {
	var b [16]byte
	var bs []byte
	bs = b[:16]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.ClientId
	bs[8] = byte(tmp32)
	bs[9] = byte(tmp32 >> 8)
	bs[10] = byte(tmp32 >> 16)
	bs[11] = byte(tmp32 >> 24)
	tmp32 = t.CmdId.SeqNum
	bs[12] = byte(tmp32)
	bs[13] = byte(tmp32 >> 8)
	bs[14] = byte(tmp32 >> 16)
	bs[15] = byte(tmp32 >> 24)
	wire.Write(bs)
	bs = b[:]
	alen1 := int64(len(t.Rep))
	if wlen := binary.PutVarint(bs, alen1); wlen >= 0 {
		wire.Write(b[0:wlen])
	}
	for i := int64(0); i < alen1; i++ {
		bs = b[:1]
		bs[0] = byte(t.Rep[i])
		wire.Write(bs)
	}
}

func (t *MReadReply) Unmarshal(rr io.Reader) error {
	var wire byteReader
	var ok bool
	if wire, ok = rr.(byteReader); !ok {
		wire = bufio.NewReader(rr)
	}
	var b [16]byte
	var bs []byte
	bs = b[:16]
	if _, err := io.ReadAtLeast(wire, bs, 16); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	t.CmdId.ClientId = int32((uint32(bs[8]) | (uint32(bs[9]) << 8) | (uint32(bs[10]) << 16) | (uint32(bs[11]) << 24)))
	t.CmdId.SeqNum = int32((uint32(bs[12]) | (uint32(bs[13]) << 8) | (uint32(bs[14]) << 16) | (uint32(bs[15]) << 24)))
	alen1, err := binary.ReadVarint(wire)
	if err != nil {
		return err
	}
	t.Rep = make([]byte, alen1)
	for i := int64(0); i < alen1; i++ {
		bs = b[:1]
		if _, err := io.ReadAtLeast(wire, bs, 1); err != nil {
			return err
		}
		t.Rep[i] = byte(bs[0])
	}
	return nil
}

func (t *MPing) BinarySize() (nbytes int, sizeKnown bool) {
	return 8, true
}

type MPingCache struct {
	mu	sync.Mutex
	cache	[]*MPing
}

func NewMPingCache() *MPingCache {
	c := &MPingCache{}
	c.cache = make([]*MPing, 0)
	return c
}

func (p *MPingCache) Get() *MPing {
	var t *MPing
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MPing{}
	}
	return t
}
func (p *MPingCache) Put(t *MPing) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MPing) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MPing) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	return nil
}

func (t *MPingRep) BinarySize() (nbytes int, sizeKnown bool) {
	return 8, true
}

type MPingRepCache struct {
	mu	sync.Mutex
	cache	[]*MPingRep
}

func NewMPingRepCache() *MPingRepCache {
	c := &MPingRepCache{}
	c.cache = make([]*MPingRep, 0)
	return c
}

func (p *MPingRepCache) Get() *MPingRep {
	var t *MPingRep
	p.mu.Lock()
	if len(p.cache) > 0 {
		t = p.cache[len(p.cache)-1]
		p.cache = p.cache[0:(len(p.cache) - 1)]
	}
	p.mu.Unlock()
	if t == nil {
		t = &MPingRep{}
	}
	return t
}
func (p *MPingRepCache) Put(t *MPingRep) {
	p.mu.Lock()
	p.cache = append(p.cache, t)
	p.mu.Unlock()
}
func (t *MPingRep) Marshal(wire io.Writer) {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	tmp32 := t.Replica
	bs[0] = byte(tmp32)
	bs[1] = byte(tmp32 >> 8)
	bs[2] = byte(tmp32 >> 16)
	bs[3] = byte(tmp32 >> 24)
	tmp32 = t.Ballot
	bs[4] = byte(tmp32)
	bs[5] = byte(tmp32 >> 8)
	bs[6] = byte(tmp32 >> 16)
	bs[7] = byte(tmp32 >> 24)
	wire.Write(bs)
}

func (t *MPingRep) Unmarshal(wire io.Reader) error {
	var b [8]byte
	var bs []byte
	bs = b[:8]
	if _, err := io.ReadAtLeast(wire, bs, 8); err != nil {
		return err
	}
	t.Replica = int32((uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]) << 16) | (uint32(bs[3]) << 24)))
	t.Ballot = int32((uint32(bs[4]) | (uint32(bs[5]) << 8) | (uint32(bs[6]) << 16) | (uint32(bs[7]) << 24)))
	return nil
}
