package mini

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/sanyuanya/dongle/data"
	"github.com/sanyuanya/dongle/entity"
	"github.com/sanyuanya/dongle/pay/common"
)

type WeChatPayNotify struct {
	Id           string   `json:"id"`
	CreateTime   string   `json:"create_time"`
	EventType    string   `json:"event_type"`
	ResourceType string   `json:"resource_type"`
	Resource     Resource `json:"resource"`
	Summary      string   `json:"summary"`
}

type Resource struct {
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	OriginalType   string `json:"original_type"`
	Nonce          string `json:"nonce"`
}

func Notify(c fiber.Ctx) error {

	serialNo := c.Get("Wechatpay-Serial")
	signature := c.Get("Wechatpay-Signature")
	timestamp := c.Get("Wechatpay-Timestamp")
	nonce := c.Get("Wechatpay-Nonce")

	if serialNo != "1A1EAB972BD01FB2C072DD11996582D1B9F66F5A" {
		log.Printf("无效的证书序列号: %s", serialNo)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    "401",
			"message": "无效的证书序列号",
		})
	}
	body, err := io.ReadAll(bytes.NewBuffer(c.Body()))

	if err != nil {
		log.Printf("无法读取请求体: %#+v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    "422",
			"message": "无法读取请求体",
		})
	}
	message := fmt.Sprintf("%s\n%s\n%s\n", timestamp, nonce, string(body))

	env := os.Getenv("ENVIRONMENT")

	certPath := ""
	switch env {
	case "production":
		certPath = "/cert"
	default:
		certPath = "/Users/sanyuanya/hjworkspace/go_dev/dongle_new/pay/cert"
	}

	publicFilePath := fmt.Sprintf("%s/wechatpay_17BDDF6F46451DE2C953B628B76D4458B00CF054.pem", certPath)
	publicKey, err := common.ReadPublicKey(publicFilePath)
	if err != nil {
		log.Printf("无法读取公钥文件: %#+v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"code":    "500",
			"message": "无法读取公钥文件",
		})
	}

	// 验证签名
	err = VerifySignature(publicKey, signature, message)
	if err != nil {
		log.Printf("签名验证失败: %#+v", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"code":    "401",
			"message": "签名验证失败",
		})
	}

	// 解析回调通知
	var notify *WeChatPayNotify
	err = json.Unmarshal(body, &notify)
	if err != nil {
		log.Printf("无法解析回调通知: %#+v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"code":    "422",
			"message": "无法解析回调通知",
		})
	}

	log.Printf("回调通知: %#+v", notify)

	if notify.EventType == "TRANSACTION.SUCCESS" {
		// 处理支付成功通知
		// 从 notify.Resource.Ciphertext 中解密出明文数据
		apiV3Key := "4a7a83dea74494415d163055d36f5064" // 请替换为实际的 API v3 密钥
		plaintext, err := DecryptResource(apiV3Key, notify.Resource.AssociatedData, notify.Resource.Nonce, notify.Resource.Ciphertext)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    "500",
				"message": "解密失败",
			})
		}
		log.Printf("解密后的回调通知: %s", plaintext)
		// 返回处理结果

		var response *entity.DecryptResourceResponse

		err = json.Unmarshal(plaintext, &response)
		if err != nil {
			log.Printf("无法解析解密后的回调通知: %#+v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    "500",
				"message": "无法解析解密后的回调通知",
			})
		}

		log.Printf("解密后的回调通知: %v", response)

		tx, err := data.Transaction()

		if err != nil {
			log.Printf("无法创建事务: %#+v", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    "500",
				"message": "无法创建事务",
			})
		}

		err = data.UpdateOrder(tx, response)

		if err != nil {
			log.Printf("无法更新订单: %#+v", err)
			tx.Rollback()
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    "500",
				"message": "无法更新订单",
			})
		}

		err = tx.Commit()

		if err != nil {
			log.Printf("无法提交事务: %#+v", err)
		}

		return c.JSON(fiber.Map{
			"code":    "200",
			"message": "处理成功",
		})
	}

	return c.JSON(fiber.Map{
		"code":    "404",
		"message": "未知的通知类型",
	})
}

func VerifySignature(publicKey *rsa.PublicKey, signature, message string) error {
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("无法解码签名: %v", err)
	}

	hashed := sha256.Sum256([]byte(message))
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decodedSignature)
	if err != nil {
		return fmt.Errorf("签名验证失败: %v", err)
	}

	return nil
}

func DecryptResource(apiV3Key, associatedData, nonce, ciphertext string) ([]byte, error) {
	key := []byte(apiV3Key)
	ciphertextDecoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("无法解码密文: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("无法创建 AES 密码块: %v", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("无法创建 GCM 模式: %v", err)
	}

	plaintext, err := aesGCM.Open(nil, []byte(nonce), ciphertextDecoded, []byte(associatedData))
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	return plaintext, nil
}
