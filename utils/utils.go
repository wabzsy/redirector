package utils

import (
	crand "crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

func GetRandomString(length int) string {
	table := []byte("0123456789abcdef")
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, table[r.Intn(len(table))])
	}
	return string(result)
}

func GenerateSigner() (ssh.Signer, error) {
	key, err := rsa.GenerateKey(crand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return ssh.NewSignerFromKey(key)
}

func Pipe(rwc1, rwc2 io.ReadWriteCloser) error {
	errCh := make(chan error, 2)
	pipe := func(dst io.WriteCloser, src io.ReadCloser, wg *sync.WaitGroup) error {
		defer func() {
			_ = dst.Close()
			wg.Done()
		}()
		_, err := io.Copy(dst, src)
		return err
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		errCh <- pipe(rwc1, rwc2, wg)
	}()
	go func() {
		errCh <- pipe(rwc2, rwc1, wg)
	}()

	wg.Wait()

	if err := <-errCh; err == nil || err == io.EOF {
		return nil
	} else {
		return err
	}
}

func Forward(conn1, conn2 net.Conn) error {
	log.Printf("[(%s)==(%s)]<===>[(%s)==(%s)]\n",
		conn2.LocalAddr(), conn2.RemoteAddr(), conn1.LocalAddr(), conn1.RemoteAddr())
	defer func() {
		log.Printf("[(%s)==(%s)]<=X=>[(%s)==(%s)]\n",
			conn2.LocalAddr(), conn2.RemoteAddr(), conn1.LocalAddr(), conn1.RemoteAddr())
	}()
	return Pipe(conn1, conn2)
}

func IsDebug() bool {
	if os.Getenv("X_DEBUG") != "" {
		return true
	}
	return false
}

func WriteLine(conn net.Conn, data []byte) error {
	data = append(data, byte('\n'))
	_, err := conn.Write(data)
	return err
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//func ForwardSession(session1 protocol.Session, session2 protocol.Session) error {
//	for {
//		if stream, err := session1.AcceptStream(); err == nil {
//			log.Println("new stream from", stream.RemoteAddr())
//			if upstream, err := session2.OpenStream(); err == nil {
//				go log.Println(Forward(stream, upstream))
//			} else {
//				return err
//			}
//		} else {
//			return err
//		}
//	}
//}

////整形转换成字节
//func IntToBytes(n int) []byte {
//	bytesBuffer := bytes.NewBuffer([]byte{})
//	_ = binary.Write(bytesBuffer, binary.BigEndian, int32(n))
//	return bytesBuffer.Bytes()
//}
//
////字节转换成整形
//func BytesToInt(b []byte) int {
//	var x int32
//	bytesBuffer := bytes.NewBuffer(b)
//	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
//	return int(x)
//}
