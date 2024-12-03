package txnet

import (
	"testing"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xmap"
)

func init() {
	Bootstrap("example.properties")
}

func TestBootstarp(t *testing.T) {
	confPath := "example.properties"
	Bootstrap(confPath)
}

const username = "test"

func TestCreateMember(t *testing.T) {
	resp, err := CreateMember(username)
	if err != nil {
		t.Error(err)
	}
	t.Log(converter.JSON(resp))
}

func TestLogin(t *testing.T) {
	redirect := Login(username)
	t.Log(redirect)
}

func TestCreateSession(t *testing.T) {
	resp, err := CreateSession(username)
	if err != nil {
		t.Error(err)
	}
	t.Log(converter.JSON(resp))
}

func TestLogin2(t *testing.T) {
	resp, err := Login2(username)
	if err != nil {
		t.Error(err)
	}
	t.Log(converter.JSON(resp))
}

func TestTransferIN(t *testing.T) {
	remitno := time.Now().Format(`20060102150405`)
	resp, err := Transfer(username, remitno, "IN", "100")
	if err != nil {
		t.Error(err)
	}
	t.Log(remitno, converter.JSON(resp))
	// {"data":{"Code":"11100","Message":"Transfer Successful"},"result":true}
}

func TestTransferOUT(t *testing.T) {
	remitno := time.Now().Format(`20060102150405`)
	resp, err := Transfer(username, remitno, "OUT", "100")
	if err != nil {
		t.Error(err)
	}
	t.Log(remitno, converter.JSON(resp))
	// {"data":{"Code":"11100","Message":"Transfer Successful"},"result":true}
}

func TestCheckTransfer(t *testing.T) {
	remitno := "20241202154959"
	resp, err := CheckTransfer(username, remitno)
	if err != nil {
		t.Error(err)
	}
	t.Log(converter.JSON(resp))
	// {"data":{"Status":1,"TransID":"2463335","TransType":"OUT","UserName":"y76test"},"result":true}
}

func TestTransferRecord(t *testing.T) {
	resp, err := TransferRecord(username, "2024/12/2", "2024/12/3")
	if err != nil {
		t.Error(err)
	}
	t.Log(converter.JSON(resp))
	// {"data":[{"Amount":"100.0000","Balance":"100.0000","CreateTime":"2024-12-02 03:42:50","Currency":"RMB","TransID":"2463312","TransType":"IN","UserName":"y76test"},{"Amount":"-100.0000","Balance":"0.0000","CreateTime":"2024-12-02 03:43:19","Currency":"RMB","TransID":"2463315","TransType":"OUT","UserName":"y76test"}],"pagination":{"Page":1,"PageLimit":500,"TotalNumber":2,"TotalPage":1},"result":true}
}

func TestCheckUsrBalance(t *testing.T) {
	resp, err := CheckUsrBalance("", "", "")
	if err != nil {
		t.Error(err)
	}
	t.Log(converter.JSON(resp))
	// {"data":[{"Balance":0,"Currency":"RMB","LoginName":"y76test"}],"pagination":{"Page":1,"PageLimit":500,"TotalNumber":1,"TotalPage":1},"result":true}
}

func TestGameUrlBy3(t *testing.T) {
	resp, err := CreateSession(username)
	if err != nil {
		t.Error(err)
		return
	}
	sessionid := resp.Str("data/sessionid")
	if len(sessionid) < 1 {
		return
	}
	resp, err = GameUrlBy3("zh-cn", sessionid, "", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(converter.JSON(resp))
	data := resp.ArrayMapDef([]xmap.M{}, "data")
	if len(data) > 0 {
		mobile := data[0].Str("mobile")
		pc := data[0].Str("pc")
		rwd := data[0].Str("rwd")
		t.Log("mobile->", mobile)
		t.Log("pc->", pc)
		t.Log("rwd->", rwd)
	}
}

func TestWagersRecordBy3(t *testing.T) {
	resp, err := WagersRecordBy3("BetTime", "2024/12/02", "00:00:00", "23:59:59", "", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(converter.JSON(resp))
	// {"data":[{"BetAmount":"20.00","Card":"","Client":"3","Commissionable":"0.00","Currency":"RMB","ExchangeRate":"1.000000","GameCode":"82","GameType":"3001","ModifiedDate":"2024-12-02 05:11:34","Origin":"MAC","Payoff":"0.0000","Portal":"0","Result":"X","ResultType":"","RoundNo":"9-33","SerialID":"581594067","UserName":"y76test","WagerDetail":"2,1:1,20.00,0.00","WagersDate":"2024-12-02 05:11:34","WagersID":"522085416659"}],"pagination":{"Page":1,"PageLimit":500,"TotalNumber":"1","TotalPage":1},"result":true}
	// {"data":[{"BetAmount":"20.00","Card":"H.5,D.7,C.4*H.6,H.8,H.3","Client":"3","Commissionable":"20.00","Currency":"RMB","ExchangeRate":"1.000000","GameCode":"82","GameType":"3001","ModifiedDate":"2024-12-02 05:12:04","Origin":"MAC","Payoff":"20.0000","Portal":"0","Result":"W","ResultType":"6,7","RoundNo":"9-33","SerialID":"581594067","UserName":"y76test","WagerDetail":"2,1:1,20.00,20.00","WagersDate":"2024-12-02 05:11:34","WagersID":"522085416659"}],"pagination":{"Page":1,"PageLimit":500,"TotalNumber":"1","TotalPage":1},"result":true}
}

func TestGetWagersSubDetailUrlBy3(t *testing.T) {
	resp, err := GetWagersSubDetailUrlBy3("522085416659", "zh-cn", username, "3001")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(converter.JSON(resp))
	t.Log(resp.Str("data/Url"))
}
