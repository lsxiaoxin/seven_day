package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser     //一个网络连接
	buf  *bufio.Writer          //带缓冲的，提高性能
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)   //再写到conn
	return &GobCodec{
		conn: conn,
		buf: buf,
		dec: gob.NewDecoder(conn),  //从conn读
		enc: gob.NewEncoder(buf),   //先写到buf
	}
}

//从conn中读取gob数据,并decode到h
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

//从conn中读取gob数据,并decode到body
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

//提供h 和 body， encode 后发送到buff， 最后 buff刷新, 发送到conn
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()

	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}

	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}

	return nil
}

//关闭 conn
func (c *GobCodec) Close() error {
	return c.conn.Close()
}