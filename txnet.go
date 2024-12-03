package txnet

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/codingeasygo/util/xmap"
	"github.com/codingeasygo/util/xprop"
	"github.com/wfunc/go/xlog"
	"github.com/wfunc/util/xhash"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

// GenerateRandomString generates a random string of the specified length.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}

	return string(result), nil
}

var (
	Verbose   = false
	shared    *xhttp.Client
	apiHost   = ""
	website   = ""
	uppername = ""
	APIM      = map[string]xmap.M{}
)

func Bootstrap(confPath string) {
	conf := xprop.NewConfig()
	conf.Load(confPath)
	if conf.Int("txsrv/verbose") == 1 {
		Verbose = true
	}
	if Verbose {
		conf.Print()
	}
	proxyAddr := conf.Str("txsrv/proxy_addr")
	timeout := conf.Int("txsrv/timeout")
	InitTxNetwork(proxyAddr, timeout)
	apiHost = conf.Str("txsrv/api_host")
	website = conf.Str("txsrv/website")
	uppername = conf.Str("txsrv/uppername")
	conf.Range("txapi", func(k string, v interface{}) {
		m := xmap.M{}
		switch reflect.ValueOf(v).Kind() {
		case reflect.String:
			_, err := converter.UnmarshalJSON(bytes.NewBufferString(v.(string)), &m)
			if err == nil {
				APIM[k] = m
			}
		}
	})
}

func InitTxNetwork(proxyAddr string, timeout int) {
	httpClient := &http.Client{}
	if len(proxyAddr) > 0 {
		proxy, err := url.Parse(proxyAddr)
		if err != nil {
			panic(err)
		}
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}
	if timeout < 1 {
		timeout = 5
	}
	httpClient.Timeout = time.Duration(timeout) * time.Second
	shared = xhttp.NewClient(httpClient)
}

func ReadKeyABC(name string) (keyA, keyB, keyC string) {
	keyABC := APIM[name]
	keyA, _ = GenerateRandomString(keyABC.Int("keyA"))
	keyB = keyABC.Str("keyB")
	keyC, _ = GenerateRandomString(keyABC.Int("keyC"))
	return
}

func CallJSON(apiName string, request xmap.M) (resp xmap.M, err error) {
	keyA, keyB, keyC := ReadKeyABC(apiName)
	if Verbose {
		xlog.Infof("%v keyA=%v keyB=%v keyC=%v", apiName, keyA, keyB, keyC)
	}
	yyyyMMDD := time.Now().Format(`20060102`)
	username := request.Str("username")
	beforeMD5 := website + username + keyB + yyyyMMDD
	if apiName == "Transfer" {
		beforeMD5 = website + username + request.Str("remitno") + keyB + yyyyMMDD
	}
	if apiName == "CheckTransfer" || apiName == "WagersRecordBy3" || apiName == "GetWagersSubDetailUrlBy3" {
		beforeMD5 = website + keyB + yyyyMMDD
	}

	if Verbose {
		xlog.Infof("%v beforeMD5=%v", apiName, beforeMD5)
	}
	key := keyA + xhash.MD5([]byte(beforeMD5)) + keyC

	request["website"] = website
	request["uppername"] = uppername
	request["key"] = key
	var reqs string
	for k := range request {
		reqs += k + "=" + request.Str(k) + "&"
	}
	reqs = reqs[:len(reqs)-1]
	resp, err = shared.GetMap(apiHost+apiName+"?%v", reqs)
	if Verbose {
		xlog.Infof("%v %v", apiName, reqs)
	}
	return
}

func CallURL(apiName string, request xmap.M) (result string) {
	keyA, keyC, keyB := ReadKeyABC(apiName)
	yyyyMMDD := time.Now().Format(`20060102`)
	username := request.Str("username")
	key := keyA + xhash.MD5([]byte(website+username+keyB+yyyyMMDD)) + keyC
	request["website"] = website
	request["uppername"] = uppername
	request["key"] = key
	var reqs string
	for k := range request {
		reqs += k + "=" + request.Str(k) + "&"
	}
	reqs = reqs[:len(reqs)-1]
	result = apiHost + apiName + "?" + reqs
	return
}

// username	是	String	会员帐号(请输入4-20个字元, 仅可输入英文字母以及数字的组合)
// ingress	否	Integer	登入来源，请填入代码(1：网页版、2：手机网页版、4：App)，预设为9：其他
func CreateMember(username string, extra ...xmap.M) (resp xmap.M, err error) {
	req := xmap.M{}
	if len(extra) > 0 {
		req = extra[0]
	}
	req["username"] = username
	return CallJSON("CreateMember", req)
}

// lang	否	String	语系：zh-cn(简中);zh-tw(繁中);en-us(英文);euc-jp(日文);ko(韩文);th(泰文) ;vi(越南文)
// page_site	否	String	视讯：live、机率：game、彩票：Ltlottery、New BB体育：nball、BB捕鱼达人、BB捕鱼大师：fisharea，若为空白则导入视讯大厅
// maintenance_page	否	Integer	0:维护时回传讯息、1:维护时导入整合页(预设为0)
// ingress	否	Integer	登入来源，请填入代码(1：网页版、2：手机网页版、4：App)，预设为9：其他
// ip	否	String	IP
func Login(username string, extra ...xmap.M) (redirect string) {
	req := xmap.M{}
	if len(extra) > 0 {
		req = extra[0]
	}
	req["username"] = username
	return CallURL("Login", req)
}

// lang	否	String	语系：zh-cn(简中);zh-tw(繁中);en-us(英文);euc-jp(日文);ko(韩文);th(泰文) ;vi(越南文)
// ingress	否	Integer	登入来源，请填入代码(1：网页版、2：手机网页版、4：App)，预设为9：其他
// ip	否	String	IP
func CreateSession(username string, extra ...xmap.M) (resp xmap.M, err error) {
	req := xmap.M{}
	if len(extra) > 0 {
		req = extra[0]
	}
	req["username"] = username
	return CallJSON("CreateSession", req)
}

// lang	否	String	语系：zh-cn(简中);zh-tw(繁中);en-us(英文);euc-jp(日文);ko(韩文);th(泰文) ;vi(越南文)
// ingress	否	Integer	登入来源，请填入代码(1：网页版、2：手机网页版、4：App)，预设为9：其他
// ip	否	String	IP
func Login2(username string, extra ...xmap.M) (resp xmap.M, err error) {
	req := xmap.M{}
	if len(extra) > 0 {
		req = extra[0]
	}
	req["username"] = username
	return CallJSON("Login2", req)
}

func Logout(username string) (resp xmap.M, err error) {
	req := xmap.M{}
	req["username"] = username
	return CallJSON("Logout", req)
}

// remitno	是	Integer	转帐序号(唯一值)，可用贵公司转帐纪录的流水号，以避免重覆转帐 < 请用int(19) ( 1~9223372036854775806)来做设定 >，别名transid
// action	是	String	IN(转入额度) OUT(转出额度)
// remit	是	Numeric	转帐额度(正数，支援至小数点后四位)
func Transfer(username, remitno, action, remit string) (resp xmap.M, err error) {
	req := xmap.M{}
	req["username"] = username
	req["remitno"] = remitno
	req["action"] = action
	req["remit"] = remit
	return CallJSON("Transfer", req)
}

func CheckTransfer(username, remitno string) (resp xmap.M, err error) {
	req := xmap.M{}
	req["username"] = username
	req["transid"] = remitno
	return CallJSON("CheckTransfer", req)
}

// transid	否	Integer	转帐序号，对应Transfer API中的remitno< 请用int(19)( 1~9223372036854775806)来做设定 >，仅能输入2年内转帐序号
// transtype	否	String	IN转入;OUT转出
// date_start	是	Datetime	开始日期ex：2012/03/21、2012-03-21
// date_end	是	Datetime	结束日期ex：2012/03/21、2012-03-21
// start_hhmmss	否	Datetime	开始时间ex：00:00:00
// end_hhmmss	否	Datetime	结束时间ex：23:59:59
// page	否	Integer	查询页数
// pagelimit	否	Integer	每页数量 查询资料时最大每页笔数全面限制为「500」 笔
func TransferRecord(username, dateStart, dateEnd string, extra ...xmap.M) (resp xmap.M, err error) {
	req := xmap.M{}
	if len(extra) > 0 {
		req = extra[0]
	}
	req["username"] = username
	req["date_start"] = dateStart
	req["date_end"] = dateEnd
	return CallJSON("TransferRecord", req)
}

// username	否	String	会员帐号
// page	否	Integer	查询页数
// pagelimit	否	Integer	每页数量
// 查询资料时最大每页笔数全面限制为「500」 笔
func CheckUsrBalance(username, page, pagelimit string) (resp xmap.M, err error) {
	req := xmap.M{}
	if len(username) > 0 {
		req["username"] = username
	}
	if len(page) > 0 {
		req["page"] = page
	}
	if len(pagelimit) > 0 {
		req["pagelimit"] = pagelimit
	}
	return CallJSON("CheckUsrBalance", req)
}

// lang	是	String	语系：zh-cn(简中);zh-tw(繁中);en-us(英文);th(泰文) ;vi(越南文)
// sessionid	是	String	会员的session ID
// gametype	否	Integer	请详查附件二
// gamecode	否	Integer	gametype有带入时，gamecode为必填
// tag	否	String	国际厅：global、区块链：blockchain
func GameUrlBy3(lang, sessionid, gametype, gamecode, tag string) (resp xmap.M, err error) {
	req := xmap.M{}
	req["lang"] = lang
	req["sessionid"] = sessionid
	if len(gametype) > 0 {
		req["gametype"] = gametype
	}
	if len(gamecode) > 0 {
		req["gamecode"] = gamecode
	}
	if len(tag) > 0 {
		req["tag"] = tag
	}
	return CallJSON("GameUrlBy3", req)
}

// action	是	String	BetTime / ModifiedTime：须选一。
// （BetTime：使用下注时间查询信息/ ModifiedTime：使用异动时间查询资讯）
// uppername	否	String	上层帐号(action=BetTime时，需强制带入)
// date	是	Datetime	日期ex：2012/03/21、2012-03-21
// action=ModifiedTime时，日期无法带入大于7天前
// starttime	是	Datetime	开始时间ex：00:00:00
// endtime	是	Datetime	结束时间ex：23:59:59
// gametype	否	Integer	请详查附件二
// page	否	Integer	查询页数
// pagelimit	否	Integer	每页数量
// 每页笔数预设为「500」笔；若 action = ModifiedTime，并选择格式为JSON，查询资料时最大每页笔数全面限制为「10000」笔
func WagersRecordBy3(action, date, starttime, endtime, gametype, page, pagelimit string) (resp xmap.M, err error) {
	req := xmap.M{}
	req["action"] = action
	req["date"] = date
	req["starttime"] = starttime
	req["endtime"] = endtime
	if len(gametype) > 0 {
		req["gametype"] = gametype
	}
	if len(page) > 0 {
		req["page"] = page
	}
	if len(pagelimit) > 0 {
		req["pagelimit"] = pagelimit
	}
	return CallJSON("WagersRecordBy3", req)
	// UserName	帐号
	// WagersID	注单号码
	// WagersDate	下注时间
	// SerialID	局号
	// RoundNo	场次
	// GameType	游戏种类
	// WagerDetail	玩法
	// GameCode	桌号
	// Result	注单结果(C:注销,X:未结算,W:赢,L:输,D:和局)
	// ResultType	开牌结果
	// Card	结果牌
	// BetAmount	下注金额
	// Payoff	派彩金额(不包含本金)
	// Currency	币别
	// ExchangeRate	与人民币的汇率
	// Commissionable	会员有效投注额
	// Origin	1-1.ios手机：MI1
	// 1-2.ios平板：MI2
	// 1-3.Android手机：MA1
	// 1-4.Android平板：MA2
	// 2.计算机下单：P
	// 3.MAC下单：MAC
	// 4.其他：O
	// ModifiedDate	注单变更时间
	// Client	开发平台(0: WEB, 1:APP, 2:Flash, 3:HTML5, 5:AIO)
	// Portal	来源入口(0:PC, 1:APP, 2:行动装置网页版, 3:AIO, 4:AIOS, 5:AIO SDK, 6:UB+PC版, 7:UB+行动装置网页版, 8:UB客制化+PC版, 9:UB客制化+行动装置网页版, 10:PWA, 11:其他)
}

// wagersid	是	Integer	注单编号
// lang	是	String	语系：zh-cn(简中);zh-tw(繁中);en-us(英文);euc-jp(日文);ko(韩文);th(泰文) ;vi(越南文)
// username	是	String	会员帐号
// gametype	是	Integer	请详查附件二
func GetWagersSubDetailUrlBy3(wagersid, lang, username, gametype string) (resp xmap.M, err error) {
	req := xmap.M{}
	req["wagersid"] = wagersid
	req["lang"] = lang
	req["username"] = username
	req["gametype"] = gametype
	return CallJSON("GetWagersSubDetailUrlBy3", req)
}
