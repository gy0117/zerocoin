package btc

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ripemd160"
)

// Version 用于生成地址的版本
const Version = byte(0x00)

// TestVersion 用于生成测试网络地址的版本 3开头
const TestVersion = byte(0x6F)

// AddressChecksumLen 用于生成地址的校验和位数
const AddressChecksumLen = 4

// P2SHVersion P2SH 类型的地址支持多重签名 m或者n开头
const P2SHVersion = byte(0x05)

// Wallet 钱包地址
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() (*Wallet, error) {
	privateKey, publicKey, err := newKeyPair()
	if err != nil {
		return nil, err
	}
	return &Wallet{privateKey, publicKey}, nil
}

// newKeyPair 通过私钥创建公钥
func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	//1.椭圆曲线算法生成私钥
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return *privateKey, nil, err
	}
	//序列化私钥
	//marshalECPrivateKey, _ := x509.MarshalECPrivateKey(privateKey)
	//fmt.Println(string(Base58Encode(marshalECPrivateKey)))

	//2.通过私钥生成公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey, nil
}

// 从公钥得到一个地址需要五步走：
// 1. 公钥经过两次哈希(SHA256+RIPEMD160)得到一个字节数组PubKeyHash
// 2. PubKeyHash+交易版本Version拼接成一个新的字节数组Version_PubKeyHash
// 3. 对Version_PubKeyHash进行两次哈希(SHA256)并按照一定规则生成校验和CheckSum
// 4. Version_PubKeyHash+CheckSum拼接成Version_PubKeyHash_CheckSum字节数组
// 5. 对Version_PubKeyHash_CheckSum进行Base58编码即可得到地址Address

// GenerateBitcoinAddress 获取钱包地址 根据公钥生成地址
func (wallet *Wallet) GenerateBitcoinAddress() []byte {
	// 1. 对公钥做两次哈希
	ripemd160HashSum := Sha256AndRipemd160Hash(wallet.PublicKey)

	// 2. 在RIPEMD-160 哈希前面加上版本字节，对于普通的比特币地址，这个字节是0x00
	versionedPayload := append([]byte{Version}, ripemd160HashSum...)

	// 3. 对版本化的负载进行双重 SHA-256 哈希，并取前4个字节作为校验和
	checkSumBytes := CheckSum(versionedPayload)

	// 4. 将校验和附加道版本化的负载后面，使用 Base58 编码 得到最终的比特币地址
	finalPayload := append(versionedPayload, checkSumBytes...)

	//5.base58编码
	return Base58Encode(finalPayload)
}

// GenerateBitcoinTestAddress 获取钱包地址 根据公钥生成地址
// 主网需要真实的钱，这里使用测试网络的地址，除了Version与主网不同，其余步骤完全相同
func (wallet *Wallet) GenerateBitcoinTestAddress() []byte {
	// 1. 对公钥做两次哈希
	ripemd160HashSum := Sha256AndRipemd160Hash(wallet.PublicKey)

	// 2.拼接版本
	versionedPayload := append([]byte{TestVersion}, ripemd160HashSum...)

	// 3. 对版本化的负载进行双重 SHA-256 哈希，并取前4个字节作为校验和
	checkSumBytes := CheckSum(versionedPayload)

	// 4. 将校验和附加道版本化的负载后面，使用 Base58 编码 得到最终的比特币地址
	finalPayload := append(versionedPayload, checkSumBytes...)

	//5.base58编码
	return Base58Encode(finalPayload)
}

func (wallet *Wallet) GenerateBitcoinPrivateKey() string {
	//序列化私钥
	marshalECPrivateKey, _ := x509.MarshalECPrivateKey(&wallet.PrivateKey)
	priBlock := pem.Block{
		Type:  "ECD PRIVATE KEY",
		Bytes: marshalECPrivateKey,
	}
	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b)
	err := pem.Encode(bw, &priBlock)
	if err != nil {
		panic(err)
	}
	bw.Flush()
	i := b.Bytes()
	return string(Base58Encode(i))
}

func (wallet *Wallet) ResetPrivateKey(key string) error {
	//反序列化私钥
	decode := Base58Decode([]byte(key))
	block, _ := pem.Decode(decode)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	wallet.PrivateKey = *privateKey
	return nil
}

// IsValidForAddress 判断地址是否有效
func (wallet *Wallet) IsValidForAddress(address []byte) bool {
	//1.base58解码地址得到版本，公钥哈希和校验位拼接的字节数组
	version_publicKey_checksumBytes := Base58Decode(address)
	//2.获取校验位和version_publicKeHash
	checkSumBytes := version_publicKey_checksumBytes[len(version_publicKey_checksumBytes)-AddressChecksumLen:]
	version_ripemd160 := version_publicKey_checksumBytes[:len(version_publicKey_checksumBytes)-AddressChecksumLen]

	//3.重新用解码后的version_ripemd160获得校验和
	checkSumBytesNew := CheckSum(version_ripemd160)

	//4.比较解码生成的校验和CheckSum重新计算的校验和
	if bytes.Compare(checkSumBytes, checkSumBytesNew) == 0 {
		return true
	}

	return false
}

// Sha256AndRipemd160Hash 将公钥进行两次哈希
func Sha256AndRipemd160Hash(publicKey []byte) []byte {
	// 1. 对公钥进行 SHA-256 哈希
	shahash256 := sha256.New()
	shahash256.Write(publicKey)
	sha256HashSum := shahash256.Sum(nil)

	// 2. 对 SHA-256 的结果进行 RIPEMD-160 哈希
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(sha256HashSum)
	return ripemd160Hash.Sum(nil)
}

// CheckSum 两次SHA-256哈希生成校验和
func CheckSum(bytes []byte) []byte {
	hash1 := sha256.Sum256(bytes)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:AddressChecksumLen]
}
