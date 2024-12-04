windows本地环境无法开启hyper-v，无法搭建docker服务，没有mysql、redis、kafaka环境，采集数据用到的消息队列与mysql存储数据全部使用文件来操作。
#采集数据消息队列选择
## 1. 应选
    kafaka
    轻量临时使用选择redis list
## 2. 替代方案
    文件
#保存数据
    mysql
#每个任务只执行一次
    使用布隆过滤器过滤
