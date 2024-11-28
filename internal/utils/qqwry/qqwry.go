package qqwry

import (
	"LinuxOnM/cmd/qqwry"
	"encoding/binary"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"strings"
)

const (
	indexLen      = 7
	redirectMode1 = 0x01
	redirectMode2 = 0x02
)

type QQwry struct {
	Data   []byte
	Offset int64
}

func NewQQwry() (*QQwry, error) {
	IpCommonDictionary := qqwry.QQwryByte
	return &QQwry{Data: IpCommonDictionary}, nil
}

// setOffset  Set the offset.
func (q *QQwry) setOffset(offset int64) {
	q.Offset = offset
}

// readData  Read data from the file.
func (q *QQwry) readData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.setOffset(offset[0])
	}
	nums := int64(num)
	end := q.Offset + nums
	dataNum := int64(len(q.Data))
	if q.Offset > dataNum {
		return nil
	}

	if end > dataNum {
		end = dataNum
	}
	rs = q.Data[q.Offset:end]
	q.Offset = end
	return
}

// searchIndex  Find the index position.
func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.readData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	for {
		mid := q.getMiddleOffset(start, end)
		buf := q.readData(indexLen, int64(mid))
		_ip := binary.LittleEndian.Uint32(buf[:4])

		if end-start == indexLen {
			offset := byteToUInt32(buf[4:])
			buf = q.readData(indexLen)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			}
			return 0
		}

		if _ip > ip {
			end = mid
		} else if _ip < ip {
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[4:])
		}
	}
}

// getMiddleOffset
func (q *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / indexLen) >> 1
	return start + records*indexLen
}

// Find IP address lookup for corresponding geolocation information
func (q *QQwry) Find(ip string) (res ResultQQwry) {
	res = ResultQQwry{}
	res.IP = ip
	if strings.Count(ip, ".") != 3 {
		return res
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return
	}

	var area []byte
	mode := q.readMode(offset + 4)
	if mode == redirectMode1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == redirectMode2 {
			c := q.readUInt24()
			area = q.readString(c)
		} else {
			area = q.readString(countryOffset)
		}
	} else if mode == redirectMode2 {
		countryOffset := q.readUInt24()
		area = q.readString(countryOffset)
	} else {
		area = q.readString(offset + 4)
	}

	enc := simplifiedchinese.GBK.NewDecoder()
	res.Area, _ = enc.String(string(area))

	return
}

// readString Get the string.
func (q *QQwry) readString(offset uint32) []byte {
	q.setOffset(int64(offset))
	data := make([]byte, 0, 30)
	for {
		buf := q.readData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

// readUInt24
func (q *QQwry) readUInt24() uint32 {
	buf := q.readData(3)
	return byteToUInt32(buf)
}

// readMode Get the offset value type.
func (q *QQwry) readMode(offset uint32) byte {
	mode := q.readData(1, int64(offset))
	return mode[0]
}

// byteToUInt32 convert byte to uint32.
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}

type ResultQQwry struct {
	IP   string `json:"ip"`
	Area string `json:"area"`
}
