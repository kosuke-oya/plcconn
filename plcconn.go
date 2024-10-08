//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock/$GOFILE
package plcconn

// TCP通信 送受信処理（netのラッパー）

import (
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

// 受信データの最大サイズ
const RESBUF_MAX_RLEN = 256

// interfaceを設定する
type IPlcConn interface {
	Connect() error
	Write(msg []byte) ([]byte, error)
	IsConnected() bool
	SetTimeOutSecond(timeOutSecond int)
	Close() error
	OpenWriteClose(msg []byte) ([]byte, error)
}

// 通信コネクションの構造体
type PlcConn struct {
	Conn          net.Conn   // 通信コネクション
	mu            sync.Mutex // 排他制御
	IpAddres      string     // 接続先IPアドレス
	Port          int        // 接続先ポート
	TimeOutSecond int        // タイムアウト時間
}

// TCPコネクションを作る最初のコマンド
// how to use
// c := tcp.NewPlcConn("192.168.31.001", 7001, 1)
func NewPlcConn(IpAddress string, port int, timeOutSecond int) IPlcConn {
	return &PlcConn{IpAddres: IpAddress, Port: port, TimeOutSecond: timeOutSecond}
}

// TCPコネクションを開くコマンド
// すでにオープンしている場合はそのまま開いた状態をキープする
func (c *PlcConn) Connect() error {
	if c == nil {
		return errors.New("client instance is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// コネクションがnilの場合はコネクションを試みます。
	// すでにコネクションが確立されている場合は、そのコネクションを使い回します。
	if err := c.connect(); err != nil {
		return err
	}
	return nil
}

// connectはコネクションがnilの場合はコネクションを試みます。
// すでにコネクションが確立されている場合は、そのコネクションを使い回します。
func (c *PlcConn) connect() error {
	// TCPコネクションを開く
	if c.Conn == nil {
		conn, err := net.DialTimeout("tcp", c.IpAddres+":"+strconv.Itoa(c.Port), time.Duration(c.TimeOutSecond)*time.Second)
		if err != nil {
			return err
		}
		c.Conn = conn
	}

	return nil
}

// # メッセージ通信
//
// 引数
//   - conn			通信するコネクション
//   - msg			送信コマンド
//
// 返り値
//   - 受信データ
//   - エラー
func (c *PlcConn) Write(msg []byte) ([]byte, error) {
	if len(msg) == 0 {
		return []byte{}, errors.New("msg is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.connect(); err != nil {
		return []byte{}, err
	}
	return c.write(msg)
}

// PLCへのwriteコマンド
func (c *PlcConn) write(msg []byte) ([]byte, error) {
	// タイムアウトを設定する
	err := c.Conn.SetDeadline(time.Now().Add(time.Duration(c.TimeOutSecond) * time.Second))
	if err != nil {
		return []byte{}, err
	}
	// コマンドを送信する
	_, errW := c.Conn.Write(msg)
	if errW != nil {
		return []byte{}, errW
	}
	// レスポンス受信
	buf := make([]byte, RESBUF_MAX_RLEN)
	_, errR := c.Conn.Read(buf)
	if errR != nil {
		return []byte{}, errR
	}
	return buf, nil
}

// # 接続確認
func (c *PlcConn) IsConnected() bool {
	if c == nil {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Conn != nil
}

// # タイムアウト時間の設定
func (c *PlcConn) SetTimeOutSecond(timeOutSecond int) {
	c.TimeOutSecond = timeOutSecond
}

// # 通信Close
//
// 引数
// なし
//
// 返り値
// なし
func (c *PlcConn) Close() error {
	if c == nil {
		return errors.New("client instance is nil")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.close()
}

// closeはコネクションがnilでない場合にコネクションをCloseします。
func (c *PlcConn) close() error {
	var err error
	if c.Conn != nil {
		err = c.Conn.Close()
		c.Conn = nil
	}
	return err
}

// # 通信OPEN～メッセージ送受信～CLOSEまで一括で実行
//
// 引数
//   - ipAddr   接続先IPアドレス
//   - port     接続先ポート
//   - msg			送信コマンド
//   - deadline	PLCから応答がない時のタイムアウトの時間
//
// 返り値
//   - 受信データ
//   - エラー
func (c *PlcConn) OpenWriteClose(msg []byte) ([]byte, error) {
	err := c.Connect()
	if err != nil {
		return []byte{}, err
	}
	defer c.Close()
	data, err := c.Write(msg)
	if err != nil {
		return []byte{}, err
	}
	return data, err
}
