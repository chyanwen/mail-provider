package db

import (
	"fmt"
	"log"
	"time"
)

func QueryAndCopyStrategys() {
	maxstrategyid, err := QueryMaxStrategyId()
	strategys, err := QueryStrategys(maxstrategyid)
	if err != nil {
		return
	}
	InsertAlarmverifications(CopyStrategys(strategys))
}

//标记已进行告警策略验证的最大告警策略id
func QueryMaxStrategyId() (int, error) {
	var maxstrategyid int
	sql := "select max(strategy_id) from alarm_verification;"
	err := DB.QueryRow(sql).Scan(&maxstrategyid)
	if err != nil {
		log.Println("ERROR:", err)
		return 0, err
	}
	return maxstrategyid, nil
}

//获取所有告警策略
func QueryStrategys(maxstrategyid int) ([]*Strategy, error) {
	strategys := []*Strategy{}
	sql := fmt.Sprintf("select * from strategy where id>%d and priority!=6;", maxstrategyid)
	rows, err := DB.Query(sql)
	if err != nil {
		log.Println("ERROR:", err)
		return strategys, err
	}
	defer rows.Close()
	for rows.Next() {
		s := Strategy{}
		err = rows.Scan(
			&s.Id,
			&s.Metric,
			&s.Tags,
			&s.MaxStep,
			&s.Priority,
			&s.Func,
			&s.Operator,
			&s.RightValue,
			&s.Note,
			&s.RunBegin,
			&s.RunEnd,
			&s.TplId,
		)

		if err != nil {
			log.Println("scan strategy err:", err)
			continue
		}

		strategys = append(strategys, &s)
	}
	return strategys, nil
}

//复制告警策略
func CopyStrategys(strategys []*Strategy) []*AlarmVerification {
	alarmverifications := []*AlarmVerification{}
	for _, strategy := range strategys {
		a := AlarmVerification{}
		if strategy.Operator == ">" {
			strategy.Operator = "<="
		} else if strategy.Operator == ">=" {
			strategy.Operator = "<="
		} else if strategy.Operator == "<" {
			strategy.Operator = ">="
		} else if strategy.Operator == "<=" {
			strategy.Operator = ">"
		} else if strategy.Operator == "==" {
			strategy.Operator = "!="
		} else {
			strategy.Operator = "=="
		}
		sql := fmt.Sprintf("insert into strategy(metric,tags,max_step,priority,func,op,right_value,note,run_begin,run_end,tpl_id) values('%s','%s',%d,6,'%s','%s','%s','告警验证使用','%s','%s',%d)",
			strategy.Metric,
			strategy.Tags,
			strategy.MaxStep,
			strategy.Func,
			strategy.Operator,
			strategy.RightValue,
			strategy.RunBegin,
			strategy.RunEnd,
			strategy.TplId,
		)
		re, err := DB.Exec(sql)
		if err != nil {
			log.Println("exec", sql, "fail", err)
			continue
		}
		insertid, _ := re.LastInsertId()
		a.Id = 0
		a.StrategyId = strategy.Id
		a.StrategyCopyId = int(insertid)
		a.CreatedTime = int(time.Now().Unix())
		a.VerificationStatus = 0
		alarmverifications = append(alarmverifications, &a)
	}
	return alarmverifications
}

//即将进行告警验证的告警策略信息入库
func InsertAlarmverifications(alarmverifications []*AlarmVerification) {
	for _, alarmverification := range alarmverifications {
		sql := fmt.Sprintf("insert into alarm_verification(strategy_id,strategy_copyid,createdtime,verification_status) values(%d,%d,%d,%d)",
			alarmverification.StrategyId,
			alarmverification.StrategyCopyId,
			alarmverification.CreatedTime,
			alarmverification.VerificationStatus,
		)
		_, err := DB.Exec(sql)
		if err != nil {
			log.Println("exec", sql, "fail", err)
			continue
		}
	}
}
