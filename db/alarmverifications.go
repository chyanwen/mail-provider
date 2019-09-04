package db

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
        "mail-provider/config"
)

//告警验证失败逻辑判断
//告警验证状态为未验证且复制告警策略生成时间超过一天没接收到告警即认为此告警策略验证失败
func AlarmVerificationsFailed() {
	alarmverifications, err := QueryAlarmverifications()
	if err != nil {
		return
	}
	for _, alarmverification := range alarmverifications {
		if alarmverification.VerificationStatus == 0 && (int(time.Now().Unix())-alarmverification.CreatedTime > config.Config().FailTimeStd) {
			UpdateStrategyStatus(alarmverification.StrategyId, 2)
		}
	}
}

//告警验证逻辑判断
func AlarmVerificationSuccess(content string) {
	var strategyids []int
	s := Strategy{}
	//获取告警内容，解析metric tags tplid
	alarmitemslice := strings.Split(content, "\r\n")
	metric := strings.Replace(strings.SplitN(alarmitemslice[3], ":", 2)[1]," ","",-1)
	tags := strings.Replace(strings.SplitN(alarmitemslice[4], ":", 2)[1]," ","",-1)
	tplid, err := strconv.Atoi(strings.SplitN(alarmitemslice[9], "/view/", 2)[1])
	//根据metric tplid获取匹配的告警策略
	sql := fmt.Sprintf("select * from strategy where metric='%s' and tpl_id=%d and priority<6", metric, tplid)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&s.Id,
			&s.Metric,
			&s.Tags,
			&s.MaxStep,
			&s.Priority,
			&s.Func,
			&s.Operator,
			&s.RightValue,
			&s.RunBegin,
			&s.RunEnd,
			&s.Note,
			&s.TplId,
		)
		if err != nil {
			log.Println("scan strategy err:", err)
			continue
		}
		//如果告警策略的tags在告警内容的tags中，即可以认为该告警内容属于此告警策略
		if strings.Contains(tags, s.Tags) {
			strategyids = append(strategyids, s.Id)
		}
	}
	for _, strategyid := range strategyids {
		UpdateStrategyStatus(strategyid, 1)
		strategycopyid, err := GetStrategyCopyId(strategyid)
		if err != nil {
			continue
		}
		DeleteStrategyCopyId(strategycopyid)
	}
}

//获取所有的告警验证记录
func QueryAlarmverifications() ([]*AlarmVerification, error) {
	alarmverifications := []*AlarmVerification{}
	sql := "select * from alarm_verification;"
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return alarmverifications, err
	}
	defer rows.Close()
	for rows.Next() {
		a := AlarmVerification{}
		err = rows.Scan(
			&a.Id,
			&a.StrategyId,
			&a.StrategyCopyId,
			&a.CreatedTime,
			&a.VerificationStatus,
		)

		if err != nil {
			log.Println("scan  alarm_verification err:", err)
			continue
		}

		alarmverifications = append(alarmverifications, &a)
	}
	return alarmverifications, nil
}

//更改告警策略的状态信息，0默认值，表示未进行验证，1表示验证通过，2表示验证失败
func UpdateStrategyStatus(strategyid int, status int) {
	sql := fmt.Sprintf("update alarm_verification set verification_status=%d where strategy_id=%d;", status, strategyid)
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}

//根据告警策略id获取复制的告警策略id
func GetStrategyCopyId(strategyid int) (int, error) {
	var strategycopyid int
	sql := fmt.Sprintf("select strategy_copyid from alarm_verification where strategy_id=%d", strategyid)
	err := DB.QueryRow(sql).Scan(&strategycopyid)
	if err != nil {
		log.Println("ERROR:", err)
		return 0, err
	}

	return strategycopyid, nil
}

//删除复制告警策略
func DeleteStrategyCopyId(strategycopyid int) {
	sql := fmt.Sprintf("delete from strategy where id=%d", strategycopyid)
	_, err := DB.Exec(sql)
	if err != nil {
		log.Println("exec", sql, "fail", err)
	}
}
