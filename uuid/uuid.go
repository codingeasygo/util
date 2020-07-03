package uuid

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"os"
	"sync/atomic"
	"time"
)

var uuidCounter uint32 = 0

//New create new uuid
func New() string {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = MachineID[0]
	b[5] = MachineID[1]
	b[6] = MachineID[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&uuidCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return hex.EncodeToString(b[:])
}

//MachineID the machine id
var MachineID = ReadMachineID()

//ReadMachineID generates and returns a machine id.
// If this function fails to get the hostname it will cause a runtime error.
func ReadMachineID() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, _ := os.Hostname()
	// if err1 != nil {
	// 	_, err2 := io.ReadFull(rand.Reader, id)
	// 	if err2 != nil {
	// 		panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
	// 	}
	// 	return id
	// }
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	return id
}

//MID is machine id
func MID() string {
	return hex.EncodeToString(MachineID)
}
