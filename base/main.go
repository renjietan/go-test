package main

import (
	"fmt"
	"unsafe"
)

// 解析固定格式的数据包
type PacketHeader struct {
	Version   uint16
	Length    uint16
	Timestamp uint32
	Checksum  uint32
}

// 使用显式类型转换
func parsePacket1(data []byte) (*PacketHeader, error) {
	// uintptr 转 int 然后再进行进行比较
	headerSize := int(unsafe.Sizeof(PacketHeader{}))
	if len(data) < headerSize {
		return nil, fmt.Errorf("数据包太小：需要 %d 字节，实际只有 %d 字节",
			headerSize, len(data))
	}

	// 零拷贝解析：直接解释内存
	header := (*PacketHeader)(unsafe.Pointer(&data[0]))
	return header, nil
}

// 使用 uintptr 比较（需要转换 len）
func parsePacket2(data []byte) (*PacketHeader, error) {
	if uintptr(len(data)) < unsafe.Sizeof(PacketHeader{}) {
		return nil, fmt.Errorf("数据包太小")
	}

	header := (*PacketHeader)(unsafe.Pointer(&data[0]))
	return header, nil
}

func main() {
	// 测试数据
	data := make([]byte, unsafe.Sizeof(PacketHeader{}))
	fmt.Println("测试数据:", data)

	// 填充测试数据
	header := (*PacketHeader)(unsafe.Pointer(&data[0]))
	header.Version = 1
	header.Length = uint16(len(data))
	header.Timestamp = 1234567890
	header.Checksum = 0xDEADBEEF

	// 测试解析
	parsed, err := parsePacket1(data)
	fmt.Println("\n测试正常数据包:", parsed)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("解析成功: Version=%d, Length=%d, Timestamp=%d, Checksum=0x%x\n",
			parsed.Version, parsed.Length, parsed.Timestamp, parsed.Checksum)
	}

	// 测试太小的数据包
	smallData := make([]byte, 4)
	_, err = parsePacket1(smallData)
	fmt.Println("\n测试小数据包:", err)
}
