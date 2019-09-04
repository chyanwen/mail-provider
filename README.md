mail-provider
==========================================================

声明：本程序是基于https://github.com/open-falcon/mail-provider修改，本人未因此产生任何收益。


友情提示：

如果使用自建邮箱，比如最简单的sendmail，新版的go环境 smtp会验证不通过，导致告警发不出，因此在编译的时候，推荐老版本的go环境，本人使用go1.9.5测试有效，另外已提供了一个编译好的二进制文件。

编译方法：

go get ./...


./contorl build

数据库表结构导入：

mysql -uusername -ppassword <alarm-verification-db-schema.sql

配置文件：

{
    "debug": true,
    "http": {
        "listen": "0.0.0.0:4000",
        "token": ""
    },
    "database": "root:@tcp(127.0.0.1:3306)/falcon_portal?loc=Local&parseTime=true", //新增配置，连接falcon db
    "maxIdle": 100, //新增配置，db使用，默认即可，抄袭于open-falcon hbs组件的配置参数
    "cronStep": 3600,  //新增配置，定时获取最新告警策略、复制告警策略、告警策略配置错误验证  
    "failTimeStd": 86400, //新增配置，用于判断告警策略配置错误的时间周期认定
    "smtp": {
        "addr": "mail.example.com:25",
        "username": "falcon@example.com",
        "password": "123456",
        "from": "falcon@example.com"
    }
}

启动：

./control start


 
