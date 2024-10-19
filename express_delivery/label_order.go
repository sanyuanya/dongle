package expressdelivery

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sanyuanya/dongle/entity"
)

func LabelOrder(param *entity.LabelOrderRequest) (*entity.LabelOrderResponse, error) {

	key := "mbzPBBLg6641"
	secret := "f969f49a93dc45979478aece402b0264"
	apiURL := "https://api.kuaidi100.com/label/order"
	method := "order"

	paramBytes, err := json.Marshal(param)
	if err != nil {
		return nil, fmt.Errorf("快递100 电子面单与云打印请求下单接口 json marshal LabelOrderRequest struct  error %v", err)
	}

	t := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	hash := md5.New()
	hash.Write([]byte(string(paramBytes) + t + key + secret))
	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))

	payload := url.Values{}
	payload.Set("key", key)
	payload.Set("method", method)
	payload.Set("t", t)
	payload.Set("param", string(paramBytes))
	payload.Set("sign", sign)

	resp, err := http.PostForm(apiURL, payload)
	if err != nil {
		return nil, fmt.Errorf("快递100 电子面单与云打印请求下单接口 http 发起请求 出错: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("快递100 电子面单与云打印下单接口 请求失败: %s", resp.Status)
	}

	var labelOrderResponse *entity.LabelOrderResponse

	if err := json.NewDecoder(resp.Body).Decode(&labelOrderResponse); err != nil {
		return nil, fmt.Errorf("快递100 电子面单与云打印请求下单接口 解析响应失败：%v", err)
	}

	return labelOrderResponse, nil
}
