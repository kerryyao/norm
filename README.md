# 组件应用标准结构定义库

#### 概述
1. 除特殊说明外，正常组件使用时在程序引导时调整用组件的Init进行配置初始化，使用New方法进行组件实例获取。

#### 组件下载工具
1. 安装
> go install myschools.me/suguo/norm/my:latest
2. 使用
> my -dl mysql

#### 组件实现
* MySQL  
  > "gorm.io/driver/mysql"  
	> "gorm.io/gorm" 
	> "gorm.io/gorm/logger"   
	> "gorm.io/plugin/dbresolver"
