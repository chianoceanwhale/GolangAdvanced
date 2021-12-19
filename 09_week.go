// 第一.总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。 
//   Socket,给程序员用的网络框架，封装了TCP和UDP传输通讯协议, 用来开发网络应用程序
//   TCP全称Transmission Control Protocol, 用于通信的协议。

// 什么是粘包
// 发送者发了两条消息：消息1:ABC, 消息2：DEF
// 正常情况，接受者收到：消息1:ABC, 消息2：DEF

// 粘包，接受者收到：消息1:ABCDEF
// 半包，接受者收到：消息1:AB, 消息2:CDEF
// 半包，接受者收到：消息1:ABCD, 消息2:EF


// 为什么会发生粘包
// TCP协议为了高效传输数据付出的代价
// 对TCP来说，它处理的是底层的数据流，数据流本身没有任何开始和结束的边界
// 发送数据过程：应用程序发送消息包，消息包以数据流的形式放入缓冲区，等缓冲区的数据流到达一定阈值后，再发送到网络上
// 接受数据过程：接受到网络过来的数据流，放入缓冲区，缓冲区的数据流到达一定阈值后，通知应用程序进行读取数据

// 在数据发送和接受的过程中，都是对数据流进行操作
// 1.在发送数据的时候
// 应用程序发送的数据长度超过缓冲区空间，这就发生数据流拆分, 同一个数据包就会通过多次发送完成，表现就是上述的半包情况
// 应用程序发送的数据小于超过缓冲区空间，等到同多个数据包填满缓冲区再进行发送，表现就是上述的粘包情况

// 2.在接受数据的时候
// 应用程序没有继续读取缓冲区的数据流，导致缓冲区放了多个数据包数据，再进行读取，也是上述的粘包

// 怎么处理粘包
// 方式1: fix length
// 发送方，每次发送固定长度的数据，并且不超过缓冲区，接受方每次按固定长度区接受数据

// 方式2: delimiter based
// 发送方，在数据包添加特殊的分隔符，用来标记数据包边界

// 方式3: length field based
// 发送方，在消息数据包头添加包长度信息

// 第二.实现一个从 socket connection 中解码出 goim 协议的解码器。 
package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	data := encoder("Hello Class!!!")
	decoder(data)
}

/*
goim 协议结构
4bytes PacketLen 包长度，在数据流传输过程中，先写入整个包的长度，方便整个包的数据读取。
2bytes HeaderLen 头长度，在处理数据时，会先解析头部，可以知道具体业务操作。
2bytes Version 协议版本号，主要用于上行和下行数据包按版本号进行解析。
4bytes Operation 业务操作码，可以按操作码进行分发数据包到具体业务当中。
4bytes Sequence 序列号，数据包的唯一标记，可以做具体业务处理，或者数据包去重。
PacketLen-HeaderLen Body 实际业务数据，在业务层中会进行数据解码和编码。
*/

func decoder(data []byte) {
	if len(data) <= 16 {
		fmt.Println("data len < 16.")
		return
	}

	packetLen := binary.BigEndian.Uint32(data[:4])
	fmt.Printf("packetLen:%v\n", packetLen)

	headerLen := binary.BigEndian.Uint16(data[4:6])
	fmt.Printf("headerLen:%v\n", headerLen)

	version := binary.BigEndian.Uint16(data[6:8])
	fmt.Printf("version:%v\n", version)

	operation := binary.BigEndian.Uint32(data[8:12])
	fmt.Printf("operation:%v\n", operation)

	sequence := binary.BigEndian.Uint32(data[12:16])
	fmt.Printf("sequence:%v\n", sequence)

	body := string(data[16:])
	fmt.Printf("body:%v\n", body)
}

func encoder(body string) []byte {
	headerLen := 16
	packetLen := len(body) + headerLen
	ret := make([]byte, packetLen)

	binary.BigEndian.PutUint32(ret[:4], uint32(packetLen))
	binary.BigEndian.PutUint16(ret[4:6], uint16(headerLen))

	version := 5
	binary.BigEndian.PutUint16(ret[6:8], uint16(version))
	operation := 6
	binary.BigEndian.PutUint32(ret[8:12], uint32(operation))
	sequence := 7
	binary.BigEndian.PutUint32(ret[12:16], uint32(sequence))

	byteBody := []byte(body)
	copy(ret[16:], byteBody)

	return ret
}
